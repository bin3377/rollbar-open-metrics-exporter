package rollbar_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/bin3377/rollbar-open-metrics-exporter/internal/rollbar"
	"github.com/sirupsen/logrus"
)

func assert(tb testing.TB, condition bool, msg string, v ...any) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		log.Printf("%s:%d: "+msg+"\n\n", append([]any{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		log.Printf("%s:%d: unexpected error: %s\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

func equals(tb testing.TB, exp, act any) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		log.Printf("%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	rollbar.AccountReadAccessToken = os.Getenv("ROLLBAR_ACCOUNT_READ_TOKEN")
	rollbar.AccountWriteAccessToken = os.Getenv("ROLLBAR_ACCOUNT_WRITE_TOKEN")
}

func Test_ListProjects(t *testing.T) {
	ps, err := rollbar.ListProjects()
	ok(t, err)
	for _, p := range ps {
		logrus.Printf("%v", p)
	}
}

func Test_ListProjectToken(t *testing.T) {
	ps, err := rollbar.ListProjectAccessTokens(224205)
	ok(t, err)
	for _, p := range ps {
		logrus.Printf("%v", p)
	}
}

func Test_GetOrCreateProjectReadToken(t *testing.T) {
	ps, err := rollbar.ListProjects()
	ok(t, err)
	for _, p := range ps {
		logrus.Printf("reading token of %d %s...", p.ID, p.Name)
		token, err := rollbar.GetOrCreateProjectReadToken(p.ID)
		ok(t, err)
		assert(t, token != nil, "token not nil")
	}
}

func Test_ListEnvrionments(t *testing.T) {
	ps, err := rollbar.ListProjects()
	ok(t, err)
	for _, p := range ps {
		logrus.Printf("reading token of %d %s...", p.ID, p.Name)
		token, err := rollbar.GetOrCreateProjectReadToken(p.ID)
		ok(t, err)
		envs, err := rollbar.ListEnvrionments(token.AccessToken)
		ok(t, err)
		for _, env := range envs {
			logrus.Printf("%v", env)
		}
	}
}

func Test_GetOccurrencesMetrics(t *testing.T) {
	token := os.Getenv("ROLLBAR_PROJECT_READ_TOKEN")
	metrics, err := rollbar.GetOccurrencesMetrics(token,
		rollbar.NewItemOccurrencesInput(time.Hour),
	)
	ok(t, err)
	for _, tp := range metrics.Timepoints {
		logrus.Printf("Time %s:", time.Unix(tp.Timestamp, 0))
		for _, row := range tp.MetricsRows {
			line := ""
			for _, cell := range row {
				line += fmt.Sprintf("%s: %v,", cell.Field, cell.Value)
			}
			logrus.Println(line)
		}
	}
}

func Test_GetItemOccurrences(t *testing.T) {
	token := os.Getenv("ROLLBAR_PROJECT_READ_TOKEN")
	occs, err := rollbar.GetItemOccurrences(token, time.Hour)
	ok(t, err)
	for _, occ := range occs {
		logrus.Printf("%v", occ)
	}
}

func Test_GetItemByID(t *testing.T) {
	token := os.Getenv("ROLLBAR_PROJECT_READ_TOKEN")
	occs, err := rollbar.GetItemOccurrences(token, time.Hour)
	ok(t, err)
	for _, occ := range occs {
		item, err := rollbar.GetItemByID(token, occ.ItemID)
		ok(t, err)
		equals(t, item.ID, occ.ItemID)
		logrus.Printf("%v", item)
	}
}

func Test_ListItemsWithIDs(t *testing.T) {
	token := os.Getenv("ROLLBAR_PROJECT_READ_TOKEN")
	occs, err := rollbar.GetItemOccurrences(token, time.Hour)
	ok(t, err)
	ids := make([]int, 0)
	for _, occ := range occs {
		ids = append(ids, occ.ItemID)
	}
	items, err := rollbar.ListItemsWithIDs(token, ids)
	ok(t, err)
	equals(t, len(occs), len(items))
	for _, item := range items {
		logrus.Printf("%v", item)
	}
}
