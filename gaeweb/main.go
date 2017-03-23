package main

import (
	"net/http"
	"strings"
	"time"

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

const kDefaultUrl = "http://www.google.com"

func create(ctx context.Context, key string) {
	entry := new(Entry)
	entry.Value = kDefaultUrl
	keyx := datastore.NewKey(ctx, "Entry", key, 0, nil)
	if _, err := datastore.Put(ctx, keyx, entry); err != nil {
		log.Warningf(ctx, "Insertion failed")
	}
}

func load(ctx context.Context, key string) string {
	var entry Entry
	keyx := datastore.NewKey(ctx, "Entry", key, 0, nil)
	if err := datastore.Get(ctx, keyx, &entry); err != nil {
		log.Warningf(ctx, "Key not found")
		// Not needed, but makes it easier to add new entries.
		// Be careful with this. A melicious user or bot could end
		// up creating a lot of datastore entries.
		create(ctx, key)
		return kDefaultUrl
	}

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

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	url := load(ctx, strings.TrimLeft(r.URL.Path, "/"))
	http.Redirect(w, r, url, 302)
}

func init() {
	http.HandleFunc("/", handler)
}
