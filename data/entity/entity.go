package entity

import "time"

type Entry struct {
	Value string `datastore:"value,noindex"`
}

type LogEntry struct {
	Key       string    `datastore:",noindex"`
	Url       string    `datastore:",noindex"`
	Timestamp time.Time `datastore:",noindex"`
}
