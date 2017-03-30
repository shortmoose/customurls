package client

import (
	"time"

	"github.com/nthnca/customurls/data/entity"

	"github.com/nthnca/datastore"
)

func LoadEntry(c datastore.Client, key string) (*entity.Entry, error) {
	var entry entity.Entry
	keyx := c.NameKey("Entry", key, nil)
	if err := c.Get(keyx, &entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

func CreateEntry(c datastore.Client, key, url string) error {
	entry := new(entity.Entry)
	entry.Value = url
	keyx := c.NameKey("Entry", key, nil)
	if _, err := c.Put(keyx, entry); err != nil {
		return err
	}
	return nil
}

func CreateLogEntry(c datastore.Client, key, url string) error {
	entry := new(entity.LogEntry)
	entry.Key = key
	entry.Url = url
	entry.Timestamp = time.Now()
	keyx := c.IncompleteKey("LogEntry", nil)
	if _, err := c.Put(keyx, entry); err != nil {
		return err
	}
	return nil
}
