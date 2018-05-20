package code

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/nthnca/customurls/config"
	"github.com/nthnca/customurls/data/client"
	"github.com/nthnca/customurls/data/entity"

	"github.com/nthnca/datastore"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type usage struct {
	key     string
	url     string
	week    int
	month   int
	allTime int
}

func showStats(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	clt := datastore.NewGaeClient(ctx)

	data := make(map[string]usage)

	var entries []entity.Entry
	keys, _ := clt.GetAll(clt.NewQuery("Entry"), &entries)

	for i := range entries {
		data[keys[i].GetName()] = usage{
			key: keys[i].GetName(),
			url: entries[i].Value}
	}

	var logs []entity.LogEntry
	_, err := clt.GetAll(clt.NewQuery("LogEntry"), &logs)
	if err != nil {
		log.Warningf(ctx, "Unable to get log entries: %v\n", err)
		return
	}

	now := time.Now()
	week := now.Add(-time.Hour * 24 * 7)
	month := now.Add(-time.Hour * 24 * 7 * 28)

	for _, log := range logs {
		e, ok := data[log.Key]
		if !ok {
			continue
		}
		if log.Timestamp.After(week) {
			e.week++
		}
		if log.Timestamp.After(month) {
			e.month++
		}
		e.allTime++
		data[log.Key] = e
	}

	var arr []usage
	for _, value := range data {
		arr = append(arr, value)
	}

	sort.Slice(arr, func(i, j int) bool {
		if arr[i].week != arr[j].week {
			return arr[i].week < arr[j].week
		}
		if arr[i].month != arr[j].month {
			return arr[i].month < arr[j].month
		}
		return arr[i].allTime < arr[j].allTime
	})

	for _, y := range arr {
		fmt.Fprintf(w, "%-15s %4d %4d %4d\n", y.key, y.week, y.month,
			y.allTime)
	}
}

// Returns a lowercase version of the given string, returns an empty string if invalid.
func getKey(str string) string {
	if !regexp.MustCompile(`^[A-Za-z0-9-]+$`).MatchString(str) {
		return ""
	}

	return strings.ToLower(str)
}

func redirectAndLog(w http.ResponseWriter, r *http.Request, key string) {
	ctx := appengine.NewContext(r)

	key = getKey(key)
	if key == "" {
		log.Warningf(ctx, "Invalid key attempted.")
		http.Redirect(w, r, config.DefaultURL, 302)
		return
	}

	clt := datastore.NewGaeClient(ctx)
	entry, err := client.LoadEntry(clt, key)
	if err != nil {
		log.Warningf(ctx, "Unable to load '%s': %v", key, err)
		return
	}

	log.Infof(ctx, "Redirecting %s to %s", key, entry.Value)
	client.CreateLogEntry(clt, key, entry.Value)
	http.Redirect(w, r, entry.Value, 302)
}

func showUrlModificationForm(w http.ResponseWriter, r *http.Request, key string) {
	ctx := appengine.NewContext(r)
	key = getKey(key)
	if key == "" {
		log.Warningf(ctx, "Invalid key attempted.")
		http.Redirect(w, r, config.DefaultURL, 302)
		return
	}

	clt := datastore.NewGaeClient(ctx)
	entry, err := client.LoadEntry(clt, key)
	if err != nil {
		log.Warningf(ctx, "Unable to load '%s': %v", key, err)
		showUrlForm(w, key, "")
		return
	}

	showUrlForm(w, key, entry.Value)
}

func showUrlForm(w http.ResponseWriter, key, url string) {
	html := `<html>
<form name="addgolink" action="/" method="post">
Key:<br />
<input type="text" name="key" value="%s"><br />
URL:<br />
<input type="text" name="url" value="%s"><br />
Validate:<br />
<input type="text" name="check"><br />
<br />
<input type="submit" value="Submit">
</form></html>`
	fmt.Fprintf(w, html, key, url)
}

func saveUrl(w http.ResponseWriter, r *http.Request) {
	// config.Check of "", means readonly system.
	if config.Check == "" {
		return
	}

	if r.FormValue("check") != config.Check {
		return
	}

	ctx := appengine.NewContext(r)
	key := getKey(r.FormValue("key"))
	if key == "" {
		log.Warningf(ctx, "Invalid Key %s", key)
		return
	}

	url := r.FormValue("url")
	if len(url) < 4 || url[:4] != "http" {
		log.Warningf(ctx, "Invalid URL %s", url)
		return
	}

	clt := datastore.NewGaeClient(ctx)
	client.CreateEntry(clt, key, url)
	return
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		saveUrl(w, r)
	} else if "/x/ls" == r.URL.Path {
		showStats(w, r)
	} else if len(r.URL.Path) > 3 && "/x/" == r.URL.Path[:3] {
		showUrlModificationForm(w, r, r.URL.Path[3:])
	} else {
		redirectAndLog(w, r, r.URL.Path[1:])
	}
}

func init() {
	http.HandleFunc("/", handler)
}
