package client

import (
	"fmt"
	"time"

	"github.com/nthnca/customurls/data/entity"
	"github.com/nthnca/datastore"
)

func LoadEntry(c datastore.Client, key string) (*entity.Entry, error) {
	var entry entity.Entry
	keyx := c.NameKey("Entry", key)
	if err := c.Get(keyx, &entry); err != nil {
		return nil, fmt.Errorf("Key not found: %v", err)
	}
	return &entry, nil
}

func CreateEntry(c datastore.Client, key, url string) error {
	entry := new(entity.Entry)
	entry.Value = url
	keyx := c.NameKey("Entry", key)
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
	keyx := c.IncompleteKey("LogEntry")
	if _, err := c.Put(keyx, entry); err != nil {
		return err
	}
	return nil
}
