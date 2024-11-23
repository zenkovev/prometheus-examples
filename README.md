# Курс по DevOps: Доклад: Prometheus

```shell
docker compose up --build -d
docker compose down --rmi=all

sudo apt update
sudo apt install prometheus

sudo netstat -tulpn
# tcp6       0      0 :::9090                 :::*                    LISTEN      72026/prometheus
# tcp6       0      0 :::9100                 :::*                    LISTEN      71571/prometheus-no

sudo vim /etc/prometheus/prometheus.yml
# scrape_configs:
# ...
#   - job_name: metrics-go-example
#     static_configs:
#       - targets: ['localhost:8001']
# 
#   - job_name: metrics-py-example
#     static_configs:
#       - targets: ['localhost:8002']
sudo systemctl restart prometheus
```

Результаты трудов:
- `http://localhost:8001/`
- `http://localhost:8002/`
- `http://localhost:9090/classic/graph`

```shell
sudo vim /etc/prometheus/prometheus-alerts.yml
# groups:
#   - name: custom_rules
#     rules:
#       - alert: PostgreSQLErrors
#         expr: rate(pg_errors_count[30s]) > 0
#         for: 0m
#         labels:
#           severity: critical

sudo vim /etc/prometheus/prometheus.yml
# rule_files:
#   - "prometheus-alerts.yml

sudo systemctl restart prometheus
```

Результаты трудов:
- `http://localhost:9090/classic/alerts`
