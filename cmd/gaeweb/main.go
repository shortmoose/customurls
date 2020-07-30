package main

import (
	"log"
	"net/http"
	"os"

	"github.com/shortmoose/customurls/internal/redirects"
	"github.com/shortmoose/customurls/internal/util"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		util.Handler(w, r, redirects.UrlMapping)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	// Don't put anything important past this next line, it won't get run.
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
