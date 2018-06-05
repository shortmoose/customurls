# customurls

Custom URL service - runs on Google appengine

## Features

This is a simple, yet useful, URL shortener/beautifier. This allows you to
assign URLs to keys so a URL like, https://\<your-domain\>/\<key\>, will redirect you
to the specified URL. This can be used for the purpose of shortening URLs,
or the reason I use it is to make it easier to remember URLs. For example
https://foo.example/expenses is a lot easier to remember than some long
Google docs URL.

- https://\<domain\>/\<key\> will redirect you to the assigned URL.
- https://\<domain\>/admin/\<key\> will allow you to assign a URL to the given key.
- https://\<domain\>/admin/ls will give you information about how often each
  key has been used.

All "admin" links are only accessible to administrators of your site. This allows
you to easily configure who can modify the URL to key mapping.

## Setting up a custom search in your browser

For even simpler use in Chrome, Firefox, and possibly other browsers you can make typing your custom URLs even simpler by setting up a custom search. Now instead of having to type 'http://domain/key', now you will be able to type something like 'cu key', here is how you do it:

### Chrome

Open Chrome, go to Settings, manage search engines, add new search engine, for keyword enter something like 'cu' (the shorter the better), for URL enter your custom URLs URL, so for example 'https://examplecustomurl.appspot.com/%s'.

Now try it out by typing 'cu <key>' in your search bar.
  
### Firefox

It seems to be a little more complicated to do this on Firefox, but I did get it to work. For instructions see:

https://superuser.com/questions/7327/how-to-add-a-custom-search-engine-to-firefox

## Installation

```shell
mkdir -p go/src/github.com/nthnca
cd go/src/github.com/nthnca
git clone https://github.com/nthnca/customurls.git
cd customurls
cp config/template.go config/config.go
# vim config/config.go
go get ./...
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
