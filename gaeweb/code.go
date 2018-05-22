package code

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/nthnca/customurls/lib/config"
	"github.com/nthnca/customurls/lib/data/client"
	"github.com/nthnca/customurls/lib/data/entity"
	"github.com/nthnca/customurls/lib/util"

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

	fmt.Fprintf(w, "<pre>\n")
	for i := range arr {
		y := arr[len(arr)-i-1]
		fmt.Fprintf(w, "  %-15s %4d %4d %4d\n", y.key, y.week, y.month,
			y.allTime)
	}
	fmt.Fprintf(w, "</pre>\n")
}

func redirectAndLog(cfg *config.Instance, w http.ResponseWriter, r *http.Request, key string) {
	ctx := appengine.NewContext(r)

	key = util.GetKey(key)
	if key == "" {
		log.Warningf(ctx, "Invalid key attempted.")
		http.Redirect(w, r, cfg.DefaultURL, 302)
		return
	}

	clt := datastore.NewGaeClient(ctx)
	entry, err := client.LoadEntry(clt, key)
	if err != nil {
		log.Warningf(ctx, "Unable to load '%s': %v", key, err)
		http.Redirect(w, r, cfg.DefaultURL, 302)
		return
	}

	log.Infof(ctx, "Redirecting %s to %s", key, entry.Value)
	client.CreateLogEntry(clt, key, entry.Value)
	http.Redirect(w, r, entry.Value, 302)
}

func showUrlModificationForm(cfg *config.Instance, w http.ResponseWriter, r *http.Request, key string) {
	ctx := appengine.NewContext(r)
	key = util.GetKey(key)
	if key == "" {
		log.Warningf(ctx, "Invalid key attempted.")
		http.Redirect(w, r, cfg.DefaultURL, 302)
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

func saveUrl(cfg *config.Instance, w http.ResponseWriter, r *http.Request) {
	// config.Check of "", means readonly system.
	if cfg.Check == "" {
		return
	}

	if r.FormValue("check") != cfg.Check {
		return
	}

	ctx := appengine.NewContext(r)
	key := util.GetKey(r.FormValue("key"))
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
	http.Redirect(w, r, url, 302)
}

func handler(w http.ResponseWriter, r *http.Request) {
	cfg := config.Get()
	size := len(cfg.AdminPath) + 2
	if r.Method == "POST" {
		saveUrl(cfg, w, r)
	} else if "/"+cfg.AdminPath+"/ls" == r.URL.Path {
		showStats(w, r)
	} else if len(r.URL.Path) > size && "/"+cfg.AdminPath+"/" == r.URL.Path[:size] {
		showUrlModificationForm(cfg, w, r, r.URL.Path[size:])
	} else {
		redirectAndLog(cfg, w, r, r.URL.Path[1:])
	}
}

func init() {
	http.HandleFunc("/", handler)
}
