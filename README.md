# customurls

Custom URL service - runs on Google appengine

## Installation

- run 'go get -u github.com/nthnca/customurls/...'
- edit 'customurls/config/template.go'
- run 'customurls upload'

## Basic Use

- the steps here are based on the default values in config/template.go
- go to http://appname.appspot.com/newkey
- fill out the form with a key, URL, and validate key of "supersecret"
- now go to http://appname.appspot.com/key and you will be redirected to the
  URL you entered

## Using the command line interface to manage your URLs

- 'customurls add key url' to add a new URL
- 'customurls ls' to see all existing URLs and usage stats
- 'customurls get' to get the URL for a given key
- 'customurls rm key' to delete a URL from your app
