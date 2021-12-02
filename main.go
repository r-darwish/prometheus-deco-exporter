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

var (
	uploadSpeed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "deco_upload_speed",
		Help: "Online devices",
	}, []string{"device"})

	downloadSpeed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "deco_download_speed",
		Help: "Online devices",
	}, []string{"device"})

	errors = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "deco_errors",
		Help: "Errors",
	})

	names  *map[string]struct{}
	exists = struct{}{}
)

func main() {
	c := deco.New(os.Getenv("DECO_EXPORTER_ADDR"))
	err := c.Authenticate(os.Getenv("DECO_EXPORTER_PASSWORD"))
	if err != nil {
		log.Fatal(err.Error())
	}

	prometheus.MustRegister(uploadSpeed)
	prometheus.MustRegister(downloadSpeed)
	prometheus.MustRegister(errors)

	go func() {
		for {
			updateDevices(c)
			time.Sleep(time.Second * 10)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9919", nil))
}

func updateDevices(c *deco.Client) {
	currentNames := make(map[string]struct{})
	devices, err := c.ClientList()
	if err != nil {
		log.Println(err)
		errors.Inc()
		err := c.Authenticate(os.Getenv("DECO_EXPORTER_PASSWORD"))
		if err != nil {
			log.Println(err)
		}

		return
	}
	for _, client := range devices.Result.ClientList {
		currentNames[client.Name] = exists

		uploadSpeed.WithLabelValues(client.Name).Set(float64(client.UpSpeed))
		downloadSpeed.WithLabelValues(client.Name).Set(float64(client.DownSpeed))
	}

	if names != nil {
		for client := range *names {
			_, present := currentNames[client]
			if !present {
				uploadSpeed.DeleteLabelValues(client)
				downloadSpeed.DeleteLabelValues(client)
			}
		}
	}

	names = &currentNames
}
