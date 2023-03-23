package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/bin3377/rollbar-open-metrics-exporter/internal/rollbar"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var (
	Port                 = 8080
	MetricsPath          = "/metrics"
	HealthPath           = "/healthz"
	ScrapeInterval       = 5 * time.Minute
	MaxItemsPerProject   = 0
	IncludeProjectsRegex = regexp.MustCompile("^.*$")
	ExcludeProjectsRegex = regexp.MustCompile("^$")
)

func main() {

	if e, ok := os.LookupEnv("LOG_LEVEL"); ok {
		if l, err := logrus.ParseLevel(e); err == nil {
			logrus.Infof("Log level from $LOG_LEVEL: %s", l)
			logrus.SetLevel(l)
		}
	}

	if e, ok := os.LookupEnv("INCLUDE_PROJECTS_REGEX"); ok {
		if r, err := regexp.Compile(e); err == nil {
			logrus.Infof("$INCLUDE_PROJECTS_REGEX: %s", e)
			IncludeProjectsRegex = r
		}
	}

	if e, ok := os.LookupEnv("EXCLUDE_PROJECTS_REGEX"); ok {
		if r, err := regexp.Compile(e); err == nil {
			logrus.Infof("$EXCLUDE_PROJECTS_REGEX: %s", e)
			ExcludeProjectsRegex = r
		}
	}

	if e, ok := os.LookupEnv("PORT"); ok {
		if p, err := strconv.Atoi(e); err == nil && p > 1024 {
			Port = p
			logrus.Infof("Port from $PORT: %d", p)
		}
	}

	if e, ok := os.LookupEnv("SCRAPE_INTERVAL"); ok {
		if d, err := time.ParseDuration(e); err == nil && d >= time.Minute {
			ScrapeInterval = d
			logrus.Infof("Scrape interval from $SCRAPE_INTERVAL: %s", d)
		}
	}

	if e, ok := os.LookupEnv("MAX_ITEMS"); ok {
		if n, err := strconv.Atoi(e); err == nil && n > 0 {
			MaxItemsPerProject = n
			logrus.Infof("Max items per project from $MAX_ITEMS: %d", n)
		}
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
