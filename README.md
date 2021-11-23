# Prometheus Deco Exporter

Provider prometheus metrics from [Deco M4](https://www.tp-link.com/us/deco-mesh-wifi/product-family/deco-m4/).

## Usage

Set the environment variables `DECO_EXPORTER_ADDR` to the address of your main deco and `DECO_EXPORTER_PASSWORD` to the password and run this daemon.

Connect Prometheus to port 9919 of the machine running this daemon

## Metrics

* `deco_download_speed` - Client download speed
* `deco_upload_speed` - Client upload speed
* `deco_errors` - Errors encountered by the daemon
