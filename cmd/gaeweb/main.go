package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/nthnca/customurls/lib/config"
	"github.com/nthnca/customurls/lib/data/client"
	"github.com/nthnca/customurls/lib/data/entity"
	"github.com/nthnca/customurls/lib/util"

	"github.com/nthnca/datastore"
)

type usage struct {
	key     string
	url     string
	week    int
	month   int
	allTime int
}

func main() {
	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func showStats(w http.ResponseWriter, r *http.Request) {
	clt, err := datastore.NewCloudClient(os.Getenv("GOOGLE_CLOUD_PROJECT"))

	data := make(map[string]usage)

	var entries []entity.Entry
	keys, _ := clt.GetAll(clt.NewQuery("Entry"), &entries)

	for i := range entries {
		data[keys[i].GetName()] = usage{
			key: keys[i].GetName(),
			url: entries[i].Value}
	}

	var logs []entity.LogEntry
	_, err = clt.GetAll(clt.NewQuery("LogEntry"), &logs)
	if err != nil {
		// log.Warningf(ctx, "Unable to get log entries: %v\n", err)
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

	fmt.Fprintf(w, "<html><pre>\n")
	for i := range arr {
		y := arr[len(arr)-i-1]
		fmt.Fprintf(w, "  %-15s %4d %4d %4d\n", y.key, y.week, y.month,
			y.allTime)
	}
	fmt.Fprintf(w, "</pre>\n")
	// url, _ := user.LogoutURL(ctx, "/")
	// fmt.Fprintf(w, `<a href="%s">sign out</a></html>`, url)
}

func redirectAndLog(cfg *config.Instance, w http.ResponseWriter, r *http.Request, key string) {

	key = util.GetKey(key)
	if key == "" {
		// log.Warningf(ctx, "Invalid key attempted.")
		http.Redirect(w, r, cfg.DefaultURL, 302)
		return
	}

	clt, err := datastore.NewCloudClient(os.Getenv("GOOGLE_CLOUD_PROJECT"))
	entry, err := client.LoadEntry(clt, key)
	if err != nil {
		// log.Warningf(ctx, "Unable to load '%s': %v", key, err)
		http.Redirect(w, r, cfg.DefaultURL, 302)
		return
	}

	// log.Infof(ctx, "Redirecting %s to %s", key, entry.Value)
	client.CreateLogEntry(clt, key, entry.Value)
	http.Redirect(w, r, entry.Value, 302)
}

func showUrlModificationForm(cfg *config.Instance, w http.ResponseWriter, r *http.Request, key string) {
	key = util.GetKey(key)
	if key == "" {
		// log.Warningf(ctx, "Invalid key attempted.")
		http.Redirect(w, r, cfg.DefaultURL, 302)
		return
	}

	clt, err := datastore.NewCloudClient(os.Getenv("GOOGLE_CLOUD_PROJECT"))
	entry, err := client.LoadEntry(clt, key)
	if err != nil {
		// log.Warningf(ctx, "Unable to load '%s': %v", key, err)
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

	key := util.GetKey(r.FormValue("key"))
	if key == "" {
		//	log.Warningf(ctx, "Invalid Key %s", key)
		return
	}

	url := r.FormValue("url")
	if len(url) == 0 {
		clt, _ := datastore.NewCloudClient(os.Getenv("GOOGLE_CLOUD_PROJECT"))
		client.DeleteEntry(clt, key)
		return
	} else if len(url) < 4 || url[:4] != "http" {
		// log.Warningf(ctx, "Invalid URL %s", url)
		return
	}

	clt, _ := datastore.NewCloudClient(os.Getenv("GOOGLE_CLOUD_PROJECT"))
	client.CreateEntry(clt, key, url)
	http.Redirect(w, r, url, 302)
}

func isValidUser(cfg *config.Instance, w http.ResponseWriter, r *http.Request) bool {
	/*
		ctx := appengine.NewContext(r)
		u := user.Current(ctx)
		if u == nil {
			log.Warningf(ctx, "No user logged in.")
			url, _ := user.LoginURL(ctx, "/"+cfg.AdminPath+"/ls")
			fmt.Fprintf(w, `<a href="%s">Sign in or register</a>`, url)
			return false
		}
		if !u.Admin {
			log.Warningf(ctx, "Not an administrator.")
			return false
		}
		log.Infof(ctx, "Administrator.")
	*/
	// return true
	return false
}

func handler(w http.ResponseWriter, r *http.Request) {
	cfg := config.Get()
	size := len(cfg.AdminPath) + 2
	if r.Method == "POST" {
		if isValidUser(cfg, w, r) {
			saveUrl(cfg, w, r)
		}
	} else if "/"+cfg.AdminPath+"/ls" == r.URL.Path {
		if isValidUser(cfg, w, r) {
			showStats(w, r)
		}
	} else if len(r.URL.Path) > size && "/"+cfg.AdminPath+"/" == r.URL.Path[:size] {
		if isValidUser(cfg, w, r) {
			showUrlModificationForm(cfg, w, r, r.URL.Path[size:])
		}
	} else {
		redirectAndLog(cfg, w, r, r.URL.Path[1:])
	}
}
