package main

import (
	"fmt"
	"time"

	"github.com/bin3377/rollbar-open-metrics-exporter/internal/rollbar"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var (
	occurrences = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "item_total_occurrences",
		Help: "This is the counter of total occurrences of an item",
	},
		[]string{
			"project_id",
			"item_id",
		},
	)

	itemStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "item_status",
		Help: "This is the status of item, value is always 1",
	}, []string{
		"id",
		"title",
		"project_id",
		"counter_id",
		"environment",
		"platform",
		"framework",
		"hash",
		"status",
		"level",
	})

	projectStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "project_status",
		Help: "This is the status of project, value is always 1",
	}, []string{
		"id",
		"name",
		"account_id",
		"status",
	})

	occurenceHistorigram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "item_occurrences",
		Help: "This is the histogram of item occurences",
	}, []string{
		"id",
		"title",
		"environment",
		"status",
		"level",
		"project",
	})
)

func startScrape() {

	prometheus.MustRegister(occurrences)
	prometheus.MustRegister(itemStatus)
	prometheus.MustRegister(projectStatus)
	prometheus.MustRegister(occurenceHistorigram)

	logrus.Infof("Start scraping with interval %s...", ScrapeInterval)

	s := func(t time.Time) {
		logrus.Infof("scraping at %s", t)
		if err := scrape(); err != nil {
			logrus.Errorf("scrape failed - %v", err)
		}
		logrus.Infof("scraping done (%s).", time.Since(t))
	}

	go func() {
		s(time.Now())
		for now := range time.Tick(ScrapeInterval) {
			s(now)
		}
	}()
}

var tokens = make(map[int]string)
var itemCounter = make(map[int]int64)

func scrape() error {
	ps, err := rollbar.ListProjects()
	if err != nil {
		logrus.Errorf("ListProjects failed - %v", err)
		return err
	}

	for _, p := range ps {

		// set project_status
		projectStatus.WithLabelValues(
			fmt.Sprintf("%d", p.ID),        /* id */
			p.Name,                         /* name */
			fmt.Sprintf("%d", p.AccountID), /* account_id */
			string(p.Status),               /* status */
		).Set(1)

		if p.Status == rollbar.StatusDisabled {
			continue
		}

		token, ok := tokens[p.ID]
		if !ok {
			t, err := rollbar.GetOrCreateProjectReadToken(p.ID)
			if err != nil {
				logrus.Errorf("GetOrCreateProjectReadToken failed - project: [%d]%s, %v", p.ID, p.Name, err)
				continue
			}
			token = t.AccessToken
			tokens[p.ID] = token
		}

		occs, err := rollbar.GetItemOccurrences(token, ScrapeInterval, MaxItemsPerProject)
		if err != nil {
			logrus.Errorf("GetItemOccurrences failed - project: [%d]%s, %v", p.ID, p.Name, err)
			delete(tokens, p.ID)
			continue
		}

		ids := make([]int, 0)
		for _, occ := range occs {
			ids = append(ids, occ.ItemID)
			occurenceHistorigram.WithLabelValues(
				fmt.Sprintf("%d", occ.ItemID), /* id */
				occ.ItemTitle,                 /* title */
				occ.Environment,               /* environment */
				occ.ItemStatus,                /* status */
				occ.ItemLevel,                 /* level */
				p.Name,                        /* project */
			).Observe(float64(occ.OccurrenceCount))
		}

		items, err := rollbar.ListItemsWithIDs(token, ids)
		if err != nil {
			logrus.Errorf("ListItemsWithIDs failed - project: [%d]%s, %v", p.ID, p.Name, err)
			delete(tokens, p.ID)
			continue
		}

		for _, item := range items {
			// set item_status
			itemStatus.WithLabelValues(
				fmt.Sprintf("%d", item.ID),        /* id */
				item.Title,                        /* title */
				fmt.Sprintf("%d", item.ProjectID), /* project_id */
				fmt.Sprintf("%d", item.CounterID), /* counter_id */
				item.Environment,                  /* environment */
				item.Platform,                     /* platform */
				item.Framework,                    /* framework */
				item.Hash,                         /* hash */
				item.Status,                       /* status */
				item.Level,                        /* level */
			).Set(1)

			old := itemCounter[item.ID]
			occurrences.WithLabelValues(
				fmt.Sprintf("%d", p.ID),    /* project_id */
				fmt.Sprintf("%d", item.ID), /* item_id */
			).Add(float64(item.TotalOccurrences - old))
			itemCounter[item.ID] = item.TotalOccurrences
		}
	}

	return nil
}
