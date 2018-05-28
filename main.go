package main

import (
	"flag"
	"fmt"
	"github.com/ThisisYang/gophercises/quiet_hn/hn"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl))

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var client hn.Client
		ids, err := client.TopItems()
		if err != nil {
			http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
			return
		}
		var stories []item
		respChan := make(chan *hn.ItemChan)
		// done channel to close the goroutine
		done := make(chan struct{})
		var wg sync.WaitGroup
		defer close(respChan)
		for _, id := range ids {
			wg.Add(1)
			go client.GetItemByChan(id, respChan, done, &wg)
		}

		for resp := range respChan {
			if resp.Err != nil {
				continue
			}
			item := parseHNItem(resp.Item)
			if isStoryLink(item) {
				stories = append(stories, item)
				if len(stories) >= numStories {
					// case we received enough resp, close done
					close(done)
					break
				}
			}
		}

		// if put wg.Wait() here, it will be slower
		// wg.Wait()

		data := templateData{
			Stories: stories,
			Time:    time.Now().Sub(start),
		}
		err = tpl.Execute(w, data)

		// wait for all goroutine before response
		wg.Wait()
		fmt.Println("all goroutine are done")

		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}
