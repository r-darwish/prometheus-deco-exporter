package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mrmarble/deco"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	c := deco.New(os.Getenv("DECO_EXPORTER_ADDR"))
	err := c.Authenticate(os.Getenv("DECO_EXPORTER_PASSWORD"))
	if err != nil {
		log.Fatal(err.Error())
	}

	var (
		deviceOnline = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "deco_device_online",
			Help: "Online devices",
		}, []string{"device", "interface", "wire_type"})

		uploadSpeed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "deco_download_speed",
			Help: "Online devices",
		}, []string{"device", "interface", "wire_type"})

		downloadSpeed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "deco_upload_speed",
			Help: "Online devices",
		}, []string{"device", "interface", "wire_type"})

		errors = prometheus.NewCounter(prometheus.CounterOpts{
			Name: "deco_errors",
			Help: "Errors",
		})
	)

	prometheus.MustRegister(deviceOnline)
	prometheus.MustRegister(uploadSpeed)
	prometheus.MustRegister(downloadSpeed)
	prometheus.MustRegister(errors)

	go func() {
		for {
			devices, err := c.ClientList()
			if err != nil {
				log.Println(err)
				errors.Inc()
				err := c.Authenticate(os.Getenv("DECO_EXPORTER_PASSWORD"))
				if err != nil {
					log.Println(err)
				}
			} else {
				for _, client := range devices.Result.ClientList {
					clientOnline := 0
					if client.Online {
						clientOnline = 1
					}

					deviceOnline.WithLabelValues(client.Name, client.Interface, client.WireType).Set(float64(clientOnline))
					uploadSpeed.WithLabelValues(client.Name, client.Interface, client.WireType).Set(float64(client.UpSpeed))
					downloadSpeed.WithLabelValues(client.Name, client.Interface, client.WireType).Set(float64(client.DownSpeed))
				}
			}

			time.Sleep(time.Second * 10)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9919", nil))
}
