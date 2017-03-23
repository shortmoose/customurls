package client

import (
	"fmt"
	"time"

	"github.com/nthnca/customurls/data/entity"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func CreateEntry(ctx context.Context, key, url string) {
	entry := new(entity.Entry)
	entry.Value = url
	keyx := datastore.NewKey(ctx, "Entry", key, 0, nil)
	if _, err := datastore.Put(ctx, keyx, entry); err != nil {
		log.Warningf(ctx, "Insertion failed")
	}
}

func LoadEntry(ctx context.Context, key string) (*entity.Entry, error) {
	var entry entity.Entry
	keyx := datastore.NewKey(ctx, "Entry", key, 0, nil)
	if err := datastore.Get(ctx, keyx, &entry); err != nil {
		return nil, fmt.Errorf("Key not found: %v", err)
	}
	return &entry, nil
}

func CreateLogEntry(ctx context.Context, key, url string) {
	entry := new(entity.LogEntry)
	entry.Key = key
	entry.Url = url
	entry.Timestamp = time.Now()
	keyx := datastore.NewIncompleteKey(ctx, "LogEntry", nil)
	if _, err := datastore.Put(ctx, keyx, entry); err != nil {
		log.Warningf(ctx, "Log entry failed")
	}
}
