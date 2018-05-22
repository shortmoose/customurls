package client

import (
	"time"

	"github.com/nthnca/customurls/lib/data/entity"

	"github.com/nthnca/datastore"
)

// LoadEntry loads the entry for the given key.
func LoadEntry(c datastore.Client, key string) (*entity.Entry, error) {
	var entry entity.Entry
	keyx := c.NameKey("Entry", key, nil)
	if err := c.Get(keyx, &entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

// CreateEntry saves a value (url) for the given key.
func CreateEntry(c datastore.Client, key, url string) error {
	entry := new(entity.Entry)
	entry.Value = url
	keyx := c.NameKey("Entry", key, nil)
	if _, err := c.Put(keyx, entry); err != nil {
		return err
	}
	return nil
}

// DeleteEntry removes the given key.
func DeleteEntry(c datastore.Client, key string) error {
	keyx := c.NameKey("Entry", key, nil)
	if err := c.Delete(keyx); err != nil {
		return err
	}
	return nil
}

// CreateLogEntry adds a log entry that a user loaded a given key.
func CreateLogEntry(c datastore.Client, key, url string) error {
	entry := new(entity.LogEntry)
	entry.Key = key
	entry.URL = url
	entry.Timestamp = time.Now()
	keyx := c.IncompleteKey("LogEntry", nil)
	if _, err := c.Put(keyx, entry); err != nil {
		return err
	}
	return nil
}
