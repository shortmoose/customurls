package main

import (
	"os"

	"github.com/nthnca/customurls/config"
	"github.com/nthnca/easybuild"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New(
		"custom-url",
		"URL shortening service that runs in Google AppEngine")
	easybuild.RegisterCommands(app, config.Path, config.ProjectID)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
