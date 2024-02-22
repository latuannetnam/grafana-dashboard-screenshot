# grafana-dashboard-screenshot
print  grafana dashboard to PDF vv.

## Init project
- go mod init grafana-dashboard-screenshot
- go get -u github.com/chromedp/chromedp
## Test:
- go build && ./grafana-dashboard-screenshot -grafana_protocol https -grafana_host loalhost -grafana_port 3000 -grafana_api_token xxxxx -grafana_prefix grafana -grafana_dashboard_id xxxx-yyyyy -grafana_variables "from=now-6h&to=now"

