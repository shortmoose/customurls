# customurls

Custom URL service - runs on Google appengine

## Installation

```shell
mkdir -p go/src/github.com/nthnca
cd go/src/github.com/nthnca
git clone https://github.com/nthnca/customurls.git
cd customurls
cp config/template.go config/config.go
# vim config/config.go
go install ./...
```

## Basic Use

To start with you will need to login to your account using gcloud. So for
example if you are using a robot account you will need to:
- create credentials
- save them to file
- set environment variable to point to file
  - For more info see
    https://developers.google.com/identity/protocols/application-default-credentials
- run 'customurls upload'
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
