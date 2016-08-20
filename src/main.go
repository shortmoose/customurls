package main

import (
	"appengine"
	"appengine/datastore"
	"log"
	"net/http"
	"strings"
)

type Entry struct {
	Value string `datastore:"value,noindex"`
}

const kDefaultUrl = "http://www.google.com"

func create(ctx appengine.Context, key string) {
	entry := new(Entry)
	entry.Value = kDefaultUrl
	keyx := datastore.NewKey(ctx, "Entry", key, 0, nil)
	if _, err := datastore.Put(ctx, keyx, entry); err != nil {
		log.Printf("Insertion failed")
	}
}

func load(ctx appengine.Context, key string) string {
	var entry Entry
	keyx := datastore.NewKey(ctx, "Entry", key, 0, nil)
	if err := datastore.Get(ctx, keyx, &entry); err != nil {
		log.Printf("Key not found")
		create(ctx, key)
		return kDefaultUrl
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
