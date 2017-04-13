package code

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/nthnca/customurls/config"
	"github.com/nthnca/customurls/data/client"

	"github.com/nthnca/datastore"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func showForm(w http.ResponseWriter) {
	html := `<html>
<form name="addgolink" action="https://%s.appspot.com">
Key:<br />
<input type="text" name="key"><br />
URL:<br />
<input type="text" name="url"><br />
Validate:<br />
<input type="text" name="check"><br />
<br />
<input type="submit" value="Submit">
</form></html>`
	fmt.Fprintf(w, html, config.ProjectID)
}

func load(ctx context.Context, key, url string) string {
	if key == "" {
		return config.DefaultURL
	}
	clt := datastore.NewGaeClient(ctx)
	entry, err := client.LoadEntry(clt, key)
	if err != nil {
		log.Warningf(ctx, "Unable to load '%s': %v", key, err)
		if len(url) > 4 && url[:4] == "http" {
			log.Infof(ctx, "Inserting %s:%s", key, url)
			client.CreateEntry(clt, key, url)
		}
		return config.DefaultURL
	}

	log.Infof(ctx, "Redirecting %s to %s", key, entry.Value)
	client.CreateLogEntry(clt, key, entry.Value)
	return entry.Value
}

func getKey(r *http.Request) string {
	key := strings.TrimLeft(r.URL.Path, "/")
	if key != "" {
		return key
	}

	if v, ok := r.URL.Query()["key"]; ok && v[0] != "" {
		return v[0]
	}

	return ""
}

func getURL(r *http.Request) string {
	if config.Check == "" {
		return ""
	}

	if v, ok := r.URL.Query()["check"]; !ok || v[0] != config.Check {
		return ""
	}

	if v, ok := r.URL.Query()["url"]; ok {
		return v[0]
	}
	return ""
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	key := strings.ToLower(getKey(r))
	if key == config.AddPageKey {
		showForm(w)
	} else {
		url := load(ctx, strings.ToLower(getKey(r)), getURL(r))
		http.Redirect(w, r, url, 302)
	}
}

func init() {
	http.HandleFunc("/", handler)
}
