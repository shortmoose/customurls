package main

import (
	"log"
	"net/http"
	"os"

	"github.com/nthnca/customurls/internal/config"
	"github.com/nthnca/customurls/internal/data/client"
	"github.com/nthnca/customurls/internal/util"
	"github.com/nthnca/datastore"
)

var (
	cfg config.Instance
)

func main() {
	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	cfg.ProjectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	cfg.DefaultURL = os.Getenv("DEFAULT_URL")
	log.Printf("INIT: ProjectID: %s, DefaultURL: %s", cfg.ProjectID, cfg.DefaultURL)

	log.Printf("Listening on port %s", port)
	// Don't put anything important past this next line, it won't get run.
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	key := util.GetKey(r.URL.Path[1:])
	if key == "" {
		log.Printf("Empty key attempted.")
		http.Redirect(w, r, cfg.DefaultURL, 302)
		return
	}

	clt, err := datastore.NewCloudClient(cfg.ProjectID)
	entry, err := client.LoadEntry(clt, key)
	if err != nil {
		log.Printf("Unable to load '%s': %v", key, err)
		http.Redirect(w, r, cfg.DefaultURL, 302)
		return
	}

	log.Printf("Redirecting %s to %s", key, entry.Value)
	client.CreateLogEntry(clt, key, entry.Value)
	http.Redirect(w, r, entry.Value, 302)
}
