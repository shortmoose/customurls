package entity

import "time"

// Entry is used to store a URL for a given Key.
type Entry struct {
	Value string `datastore:"value,noindex"`
	Count int    `datastore:"count,noindex"`
}

// LogEntry keeps track of usage data for each key.
type LogEntry struct {
	Key       string    `datastore:",noindex"`
	URL       string    `datastore:"Url,noindex"`
	Timestamp time.Time `datastore:",noindex"`
}
