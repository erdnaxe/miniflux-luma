package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/feeds"
	"miniflux.app/client"
)

var miniflux *client.Client
var minifluxEndpoint string
var feedTitle string

func httpHandler(w http.ResponseWriter, r *http.Request) {
	// Get new entries
	entries, err := miniflux.Entries(&client.Filter{
		Limit:     10,
		Order:     "published_at",
		Direction: "desc",
		Starred:   true,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create atom feed
	feed := &feeds.Feed{
		Title:   feedTitle,
		Link:    &feeds.Link{Href: minifluxEndpoint},
		Created: time.Now(),
		Items:   []*feeds.Item{},
	}
	for _, entry := range entries.Entries {
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       entry.Title,
			Link:        &feeds.Link{Href: entry.URL},
			Description: entry.Content,
			Author:      &feeds.Author{Name: entry.Author},
			Created:     entry.Date,
		})
	}

	// Print atom feed
	atom, err := feed.ToAtom()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, atom)
}

func main() {
	APIToken := ""
	listenAddress := ""
	certFile := ""
	keyFile := ""

	// Read command line arguments
	flag.StringVar(&minifluxEndpoint, "endpoint", "https://miniflux.example.org", "Miniflux server endpoint")
	flag.StringVar(&APIToken, "api-token", "", "Miniflux API token")
	flag.StringVar(&listenAddress, "listen-addr", "127.0.0.1:8080", "Listen on this address")
	flag.StringVar(&feedTitle, "feed-title", "Starred entries", "Title of the Atom feed")
	flag.StringVar(&certFile, "tls-cert", "", "TLS certificate file path (skip to disable TLS)")
	flag.StringVar(&keyFile, "tls-key", "", "TLS key file path (skip to disable TLS)")
	flag.Parse()

	// Authentication using API token then fetch starred items
	miniflux = client.New(minifluxEndpoint, APIToken)

	// Start web server
	http.HandleFunc("/", httpHandler)
	log.Printf("Listening on %s\n", listenAddress)
	if certFile != "" && keyFile != "" {
		log.Fatal(http.ListenAndServeTLS(listenAddress, certFile, keyFile, nil))
	} else {
		log.Fatal(http.ListenAndServe(listenAddress, nil))
	}
}
