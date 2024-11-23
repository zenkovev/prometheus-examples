import psycopg2
import time
import traceback

import prometheus_client as pc

def postgres_query() -> int:
    errors_count = 0

    conn = None
    cur = None

    try:
        conn = psycopg2.connect(
            dbname='postgres', user='postgres', password='postgres',
            host='databasehost', port=5432,
            sslmode='disable'
        )
        conn.set_session(autocommit=True)
        # To avoid performing SELECT-query in transaction mode
        # Thanks to default settings
        cur = conn.cursor()

        cur.execute('SELECT 1;')
        row = cur.fetchone()
        if len(row) != 1 or row[0] != 1:
            raise ValueError('Incorrect response to request')

    except Exception:
        traceback.print_exc()
        errors_count += 1

    finally:
        if cur is not None:
            cur.close()
        if conn is not None:
            conn.close()

    return errors_count

class MetricsCollector:
    def __init__(self):
        self.err_count = pc.Counter(
            'pg_errors_count', 'Total count of errors'
        )
        self.err_count_with_label = pc.Counter(
            'pg_errors_count_with_label', 'Total count of errors', ['info']
        )

        self.request_time_gauge = pc.Gauge(
            'pg_request_time_ms_gauge', 'Request time'
        )
        self.request_time_histogram = pc.Histogram(
            'pg_request_time_ms_histogram', 'Request time', buckets=[5, 10, 15]
        )
        # summary with quantiles seems to be usable only with another library
        # can try: https://github.com/RefaceAI/prometheus-summary/
        self.request_time_summary = pc.Summary(
            'pg_request_time_ms_summary', 'Request time'
        )

    def do_one_iteration(self):
        begin = time.time() * 1000
        err_count = postgres_query()
        end = time.time() * 1000
        duration_ms = end - begin

        self.err_count.inc(err_count)
        self.err_count_with_label.labels('first').inc(err_count)
        self.err_count_with_label.labels('second').inc(err_count)

        self.request_time_gauge.set(duration_ms)
        self.request_time_histogram.observe(duration_ms)
        self.request_time_summary.observe(duration_ms)

    def run_all_iterations(self):
        while True:
            self.do_one_iteration()
            time.sleep(3)

if __name__ == '__main__':
    mc = MetricsCollector()
    pc.start_http_server(8080)
    mc.run_all_iterations()
