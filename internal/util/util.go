package util

import (
	"log"
	"net/http"
	"regexp"
	"strings"
)

// GetKey returns a lowercase version of the given string, returns an empty string if invalid.
func GetKey(str string) string {
	if !regexp.MustCompile(`^[A-Za-z0-9-]+$`).MatchString(str) {
		return ""
	}

	return strings.ToLower(str)
}

func Handler(w http.ResponseWriter, r *http.Request, urlmap map[string]string) {
	key := GetKey(r.URL.Path[1:])
	if key == "" {
		log.Printf("Empty key attempted.")
		http.NotFound(w, r)
		return
	}

	entry := urlmap[key]
	if entry == "" {
		log.Printf("Invalid key: '%s'.", key)
		http.NotFound(w, r)
		return
	}

	log.Printf("Redirecting %s to %s", key, entry)
	http.Redirect(w, r, entry, http.StatusFound)
}
