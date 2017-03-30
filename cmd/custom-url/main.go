package main

import (
	"fmt"
	"log"
	"os"

	"google.golang.org/api/iterator"

	"github.com/nthnca/customurls/config"
	"github.com/nthnca/customurls/data/client"
	"github.com/nthnca/customurls/data/entity"
	"github.com/nthnca/datastore"
	"github.com/nthnca/easybuild"
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
	easybuild.RegisterCommands(app, config.Path, config.ProjectID)

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

func ls(_ *kingpin.ParseContext) error {
	clt, err := datastore.NewCloudClient(config.ProjectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}

	for it := clt.Run(clt.NewQuery("Entry")); ; {
		var entry entity.Entry
		k, err := it.Next(&entry)
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Printf("%v %v\n", err, iterator.Done)

			return nil
		}
		fmt.Printf("%v:%v\n", k.GetName(), entry)
	}

	return nil
}
