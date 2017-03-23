package main

import (
	"github.com/nthnca/customurls/config"
	"github.com/nthnca/easybuild"
)

func main() {
	easybuild.Build(config.Path, config.ProjectID)
}
