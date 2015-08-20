// (c) Chi Vinh Le <cvl@chinet.info> – 13.06.2015

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
)

var discovery Discovery
var domain string

func serve(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("go-get") != "1" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		if r.URL.Path != "/favicon.ico" {
			log.Printf("Bad request from address %s with URI %s", r.RemoteAddr, r.RequestURI)
		}
		return
	}

	repo := r.URL.Path[1:]
	s, e := discovery.GetRepository(repo)
	if s == "" {
		if e != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Fatalln(e.Error())
		} else {
			log.Printf("Address %s requested repository %s was not found", r.RemoteAddr, repo)
			http.NotFound(w, r)
		}
		return
	}

	var data = struct {
		Domain     string
		Repository string
		Redirect   string
	}{domain, repo, s}
	respond(w, tmpl, data)
}

// respond responds to a request by executing the html template t with data.
func respond(w http.ResponseWriter, t *template.Template, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.Execute(w, data); err != nil {
		log.Print(err)
	}
}

func main() {
	var bbUserName, bbConsumerKey, bbConsumerSecret, bbAccessToken, bbAccessTokenSecret string
	if val := os.Getenv("GOO_DOMAIN"); val != "" {
		domain = val
	} else {
		log.Fatal("GOO_DOMAIN env variable not set")
	}
	if val := os.Getenv("GOO_BITBUCKET_USERNAME"); val != "" {
		bbUserName = val
	} else {
		log.Fatal("GOO_BITBUCKET_USERNAME env variable not set")
	}
	if val := os.Getenv("GOO_BITBUCKET_CONSUMER_KEY"); val != "" {
		bbConsumerKey = val
	} else {
		log.Fatal("GOO_BITBUCKET_CONSUMER_KEY env variable not set")
	}
	if val := os.Getenv("GOO_BITBUCKET_CONSUMER_SECRET"); val != "" {
		bbConsumerSecret = val
	} else {
		log.Fatal("GOO_BITBUCKET_CONSUMER_SECRET env variable not set")
	}
	if val := os.Getenv("GOO_BITBUCKET_ACCESS_TOKEN"); val != "" {
		bbAccessToken = val
	} else {
		log.Fatal("GOO_BITBUCKET_ACCESS_TOKEN env variable not set")
	}
	if val := os.Getenv("GOO_BITBUCKET_ACCESS_TOKEN_SECRET"); val != "" {
		bbAccessTokenSecret = val
	} else {
		log.Fatal("GOO_BITBUCKET_ACCESS_TOKEN_SECRET env variable not set")
	}

	discovery = NewBitbucketDiscovery(
		bbUserName,
		bbConsumerKey, bbConsumerSecret,
		bbAccessToken, bbAccessTokenSecret,
	)

	var httpAddr = flag.String("addr", "0.0.0.0:8080", "HTTP server address")
	flag.Parse()
	http.HandleFunc("/", serve)
	fmt.Printf("Listening on %s\n", *httpAddr)
	if err := http.ListenAndServe(*httpAddr, nil); err != nil {
		log.Fatalf("Error listening, %v", err)
	}
}

var (
	tmpl = template.Must(template.New("home").Parse(
		`<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
    <meta name="robots" content="noindex, nofollow" />
    <meta name="go-import" content="{{.Domain}}/{{.Repository}} git {{.Redirect}}">
</head>
<body>
    Nothing to see.
</body>
</html>`))
)
