# customurls

[![Go Report Card](https://goreportcard.com/badge/shortmoose/customurls)](https://goreportcard.com/report/shortmoose/customurls)
[![Releases](https://img.shields.io/github/release-pre/shortmoose/customurls.svg?sort=semver)](https://github.com/shortmoose/customurls/releases)
[![LICENSE](https://img.shields.io/github/license/shortmoose/customurls.svg)](https://github.com/shortmoose/customurls/blob/master/LICENSE)
[![Go](https://github.com/shortmoose/customurls/workflows/Go/badge.svg)](https://github.com/shortmoose/customurls/actions?query=workflow%3AGo)

URL shortening (custom URL) service - runs on Google App Engine.


## Features

This is a simple, yet useful, URL shortener. This allows you to
assign URLs to keys so a URL like, https://\<your-domain\>/\<key\>, will redirect you
to the specified URL. This can be used for the purpose of shortening URLs for sharing with others,
or the reason I use it is to make it easier to remember URLs. For example
https://foo.example/expenses is a lot easier to remember than some long
Google docs URL.

- https://\<domain\>/\<key\> will redirect you to the assigned URL.


## Installation

This assumes you already know how to create projects and deploy apps to appengine.

```shell
git clone https://github.com/shortmoose/customurls.git
cd customurls
gcloud app deploy --project=<project_id> cmd/gaeweb/
go install ./...  # This will install the command line tools 'customurls'
```


## Using the command line interface to add or remove key/URL mappings

You can use the customurls command to edit your key/URL mappings. It will need a robot
account with which to do so:
- create credentials
- save them to file
- set environment variable to point to file
  - For more info see
    https://developers.google.com/identity/protocols/application-default-credentials
- Then a command like this should work: `PROJECT_ID=<project_id> GOOGLE_APPLICATION_CREDENTIALS=<path_to_json> ./customurls ls`

- 'customurls add key url' to add a new URL
- 'customurls ls' to see all existing URLs and usage stats
- 'customurls get' to get the URL for a given key
- 'customurls rm key' to delete a URL from your app


## Enjoy

Hopefully at this point you can use the `customurls` command to configure the set of key/URL mappings
your system uses, and you can use https://<project>.appspot.com/<key> to use these redirects.


## Setting up a custom search in your browser

For even simpler use in Chrome, Firefox, and possibly other browsers you can make typing your custom URLs even simpler by setting up a custom search. Now instead of having to type 'http://domain/key', now you will be able to type something like 'cu key', here is how you do it:

### Chrome

Open Chrome, go to Settings, manage search engines, add new search engine, for keyword enter something like 'cu' (the shorter the better), for URL enter your custom URLs URL, so for example 'https://examplecustomurl.appspot.com/%s'.

Now try it out by typing 'cu <key>' in your search bar.
  
### Firefox

It seems to be a little more complicated to do this on Firefox, but I did get it to work. For instructions see:

https://superuser.com/questions/7327/how-to-add-a-custom-search-engine-to-firefox
