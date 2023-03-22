package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/bin3377/rollbar-open-metrics-exporter/internal/rollbar"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var (
	Port                    = 8080
	MetricsPath             = "/metrics"
	HealthPath              = "/healthz"
	ScrapeInterval          = 5 * time.Minute
	LastScrpeAt             = time.Now()
	RollbarAccountReadToken = ""
	MaxItemsPerProject      = 0
)

func main() {

	logrus.SetLevel(logrus.InfoLevel)
	envLevel := os.Getenv("LOG_LEVEL")
	if l, err := logrus.ParseLevel(envLevel); err == nil {
		logrus.Infof("Log level from $LOG_LEVEL: %s", l)
		logrus.SetLevel(l)
	}

	envPort := os.Getenv("PORT")
	if p, err := strconv.Atoi(envPort); err == nil && p > 1024 {
		Port = p
		logrus.Infof("Port from $PORT: %d", p)
	}

	envInterval := os.Getenv("SCRAPE_INTERVAL")
	if d, err := time.ParseDuration(envInterval); err == nil && d >= time.Minute {
		ScrapeInterval = d
		logrus.Infof("Scrape interval from $SCRAPE_INTERVAL: %s", d)
	}

	envMaxItems := os.Getenv("MAX_ITEMS")
	if n, err := strconv.Atoi(envMaxItems); err == nil && n > 0 {
		MaxItemsPerProject = n
		logrus.Infof("Max items per project from $MAX_ITEMS: %d", n)
	}

	rollbar.AccountReadAccessToken = os.Getenv("ROLLBAR_ACCOUNT_READ_TOKEN")
	rollbar.AccountWriteAccessToken = os.Getenv("ROLLBAR_ACCOUNT_WRITE_TOKEN")

	startScrape()
	startHandlers()
}

func startHandlers() error {

	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.Handle("/metrics", promhttp.Handler())

	strPort := fmt.Sprintf(":%d", Port)
	logrus.Infof("Start listening on %s...", strPort)
	return http.ListenAndServe(strPort, nil)
}
