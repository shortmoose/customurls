package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/nthnca/customurls/config"
)

func runCmd(command_line string) {
	args := strings.Split(command_line, " ")
	log.Printf("Running cmd: %q", args)

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error running command")
	}
}

func run() {
	runCmd("goapp serve -port 8002 -host 192.168.86.2 gaeweb/")
}

func upload() {
	runCmd(fmt.Sprintf("goapp deploy -application %s gaeweb/",
		config.ProjectID))
}

func watch() {
	for {
		runCmd("go generate ./...")
		runCmd("go test ./...")
		runCmd("go install ./...")
		runCmd("inotifywait -qr -e modify,create,delete .")
	}
}

func main() {
	var w bool = false
	flag.BoolVar(&w, "w", false, "")
	flag.BoolVar(&w, "watch", false, "Watch")

	var r bool = false
	flag.BoolVar(&r, "r", false, "")
	flag.BoolVar(&r, "run", false, "Run")

	var u bool = false
	flag.BoolVar(&u, "u", false, "")
	flag.BoolVar(&u, "upload", false, "Upload")

	flag.Parse()

	os.Chdir(config.Path)
	if w {
		watch()
	} else if r {
		run()
	} else if u {
		upload()
	}
}
