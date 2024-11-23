package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func PostgresQuery() int {
	cfg, err := pgx.ParseConfig("postgres://postgres:postgres@databasehost:5432/postgres")
	if err != nil {
		log.Println(err)
		return 1
	}
	cfg.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	conn, err := pgx.ConnectConfig(context.Background(), cfg)
	if err != nil {
		log.Println(err)
		return 1
	}
	defer conn.Close(context.Background())

	row := conn.QueryRow(context.Background(), "SELECT 1;")
	var ans int
	if err := row.Scan(&ans); err != nil {
		log.Println(err)
		return 1
	}
	if ans != 1 {
		log.Println("incorrect response to request")
		return 1
	}

	return 0
}

type MetricsCollector struct {
	ErrCount          prometheus.Counter
	ErrCountWithLabel *prometheus.CounterVec

	RequestTimeGauge     prometheus.Gauge
	RequestTimeSummary   prometheus.Summary
	RequestTimeHistogram prometheus.Histogram
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		ErrCount: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "pg_errors_count",
				Help: "Total count of errors",
			},
		),
		ErrCountWithLabel: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "pg_errors_count_with_label",
				Help: "Total count of errors",
			},
			[]string{"info"},
		),

		RequestTimeGauge: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "pg_request_time_ms_gauge",
				Help: "Request time",
			},
		),
		RequestTimeHistogram: promauto.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "pg_request_time_ms_histogram",
				Help:    "Request time",
				Buckets: []float64{5, 10, 15},
			},
		),
		RequestTimeSummary: promauto.NewSummary(
			prometheus.SummaryOpts{
				Name: "pg_request_time_ms_summary",
				Help: "Request time",
				Objectives: map[float64]float64{
					0.1: 0.01, 0.5: 0.01, 0.9: 0.01,
				},
			},
		),
	}
}

func (mc *MetricsCollector) DoOneIteration() {
	begin := time.Now()
	errCount := PostgresQuery()
	end := time.Now()
	durationMs := float64(end.Sub(begin).Microseconds()) / 1000

	mc.ErrCount.Add(float64(errCount))
	mc.ErrCountWithLabel.With(prometheus.Labels{"info": "first"}).Add(float64(errCount))
	mc.ErrCountWithLabel.With(prometheus.Labels{"info": "second"}).Add(float64(errCount))

	mc.RequestTimeGauge.Set(durationMs)
	mc.RequestTimeSummary.Observe(durationMs)
	mc.RequestTimeHistogram.Observe(durationMs)
}

func (mc *MetricsCollector) RunAllIterations() {
	for {
		mc.DoOneIteration()
		time.Sleep(3 * time.Second)
	}
}

func main() {
	mc := NewMetricsCollector()
	go mc.RunAllIterations()

	http.Handle("/", promhttp.Handler())
	http.ListenAndServe("0.0.0.0:8080", nil)
}
