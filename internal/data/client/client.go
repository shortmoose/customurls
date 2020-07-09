package client

import (
	"context"
	"time"

	"github.com/nthnca/customurls/internal/data/entity"

	"cloud.google.com/go/datastore"
)

// LoadEntry loads the entry for the given key.
func LoadEntry(ctx context.Context, clt *datastore.Client, key string) (*entity.Entry, error) {
	var entry entity.Entry
	keyx := datastore.NameKey("Entry", key, nil)
	if err := clt.Get(ctx, keyx, &entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

// CreateEntry saves a value (url) for the given key.
func CreateEntry(ctx context.Context, clt *datastore.Client, key, url string) error {
	entry := new(entity.Entry)
	entry.Value = url
	keyx := datastore.NameKey("Entry", key, nil)
	if _, err := clt.Put(ctx, keyx, entry); err != nil {
		return err
	}
	return nil
}

/*
// DeleteEntry removes the given key.
func DeleteEntry(c datastore.Client, key string) error {
	keyx := c.NameKey("Entry", key, nil)
	if err := c.Delete(keyx); err != nil {
		return err
	}
	return nil
}
*/

// CreateLogEntry adds a log entry that a user loaded a given key.
func CreateLogEntry(ctx context.Context, clt *datastore.Client, key, url string) error {
	entry := new(entity.LogEntry)
	entry.Key = key
	entry.URL = url
	entry.Timestamp = time.Now()
	keyx := datastore.IncompleteKey("LogEntry", nil)
	if _, err := clt.Put(ctx, keyx, entry); err != nil {
		return err
	}
	return nil
}
