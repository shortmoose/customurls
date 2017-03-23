package code

import (
	"net/http"
	"strings"
	"time"

	"github.com/nthnca/customurls/config"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type Entry struct {
	Value string `datastore:"value,noindex"`
}

type LogEntry struct {
	Key       string    `datastore:",noindex"`
	Url       string    `datastore:",noindex"`
	Timestamp time.Time `datastore:",noindex"`
}

func create(ctx context.Context, key, url string) {
	entry := new(Entry)
	entry.Value = url
	keyx := datastore.NewKey(ctx, "Entry", key, 0, nil)
	if _, err := datastore.Put(ctx, keyx, entry); err != nil {
		log.Warningf(ctx, "Insertion failed")
	}
}

func load(ctx context.Context, key, url string) string {
	var entry Entry
	keyx := datastore.NewKey(ctx, "Entry", key, 0, nil)
	if err := datastore.Get(ctx, keyx, &entry); err != nil {
		log.Warningf(ctx, "Key not found: %s", key)
		if len(url) > 4 && url[:4] == "http" {
			log.Infof(ctx, "Inserting %s:%s", key, url)
			create(ctx, key, url)
		}
		return config.DefaultUrl
	}

	log.Infof(ctx, "Redirecting %s:%s", key, entry.Value)
	log_entry := new(LogEntry)
	log_entry.Key = key
	log_entry.Url = entry.Value
	log_entry.Timestamp = time.Now()
	key2 := datastore.NewIncompleteKey(ctx, "LogEntry", nil)
	if _, err := datastore.Put(ctx, key2, log_entry); err != nil {
		log.Warningf(ctx, "Log entry failed")
	}

	return entry.Value
}

func getNewUrl(r *http.Request) string {
	if config.Check == "" {
		return ""
	}

	if v, ok := r.URL.Query()["pass"]; !ok || v[0] != config.Check {
		return ""
	}

	if v, ok := r.URL.Query()["url"]; ok {
		return v[0]
	}
	return ""
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	url := load(ctx, strings.TrimLeft(r.URL.Path, "/"), getNewUrl(r))
	http.Redirect(w, r, url, 302)
}

func init() {
	http.HandleFunc("/", handler)
}
