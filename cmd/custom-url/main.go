package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/nthnca/customurls/config"
	"github.com/nthnca/customurls/data/client"
	"github.com/nthnca/customurls/data/entity"

	"github.com/nthnca/datastore"
	"github.com/nthnca/gobuild"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type keyContext struct {
	keyArg *string
	urlArg *string
}

func main() {
	app := kingpin.New(
		"custom-url",
		"URL shortening service that runs in Google AppEngine")
	gobuild.RegisterCommands(app, config.Path, config.ProjectID)

	get := &keyContext{}
	getCmd := app.Command("get", "Get URL of given key").Action(get.get)
	get.keyArg = getCmd.Arg("key", "Key").Required().String()

	set := &keyContext{}
	setCmd := app.Command("set", "Set URL of given key").Action(set.set)
	set.keyArg = setCmd.Arg("key", "Key").Required().String()
	set.urlArg = setCmd.Arg("url", "URL").Required().String()

	app.Command("ls", "list entries").Action(ls)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}

func (c *keyContext) get(_ *kingpin.ParseContext) error {
	clt, err := datastore.NewCloudClient(config.ProjectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}

	v, err := client.LoadEntry(clt, *c.keyArg)
	if err != nil {
		log.Fatalf("Unable to get entry: %v\n", err)
	}

	log.Printf("Key: %s\n", v.Value)
	return nil
}

func (c *keyContext) set(_ *kingpin.ParseContext) error {
	clt, err := datastore.NewCloudClient(config.ProjectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}

	if err := client.CreateEntry(clt, *c.keyArg, *c.urlArg); err != nil {
		log.Fatalf("Unable to set entry: %v\n", err)
	}

	log.Printf("Set %v to %v\n", *c.keyArg, *c.urlArg)
	return nil
}

type test struct {
	key     string
	url     string
	week    int
	month   int
	allTime int
}

func ls(_ *kingpin.ParseContext) error {
	clt, err := datastore.NewCloudClient(config.ProjectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}

	data := make(map[string]test)

	var entries []entity.Entry
	keys, _ := clt.GetAll(clt.NewQuery("Entry"), &entries)

	for i := range entries {
		data[keys[i].GetName()] = test{
			key: keys[i].GetName(),
			url: entries[i].Value}
	}

	var logs []entity.LogEntry
	clt.GetAll(clt.NewQuery("LogEntry"), &logs)

	now := time.Now()
	week := now.Add(-time.Hour * 24 * 7)
	month := now.Add(-time.Hour * 24 * 7 * 28)

	for _, log := range logs {
		e, ok := data[log.Key]
		if !ok {
			continue
		}
		if log.Timestamp.After(week) {
			e.week++
		}
		if log.Timestamp.After(month) {
			e.month++
		}
		e.allTime++
		data[log.Key] = e
	}

	var arr []test
	for _, value := range data {
		arr = append(arr, value)
	}

	sort.Slice(arr, func(i, j int) bool {
		if arr[i].week != arr[j].week {
			return arr[i].week < arr[j].week
		}
		if arr[i].month != arr[j].month {
			return arr[i].month < arr[j].month
		}
		return arr[i].allTime < arr[j].allTime
	})

	for _, y := range arr {
		fmt.Printf("%-15s %4d %4d %4d\n", y.key, y.week, y.month,
			y.allTime)
	}

	return nil
}
