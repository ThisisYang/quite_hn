package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/ThisisYang/gophercises/quiet_hn/hn"
	"html/template"
	"log"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()
	var client hn.Client
	pool := NewPool(numStories, &client)
	defer pool.Stop()

	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl, &client))

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(numStories int, tpl *template.Template, client *hn.Client) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		stories, err := getStories(client, numStories)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := templateData{
			Stories: stories,
			Time:    time.Now().Sub(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

func getStories(client *hn.Client, numStories int) ([]item, error) {
	ids, err := client.TopItems()
	if err != nil {
		return nil, errors.New("Failed to load top stories")
	}
	var stories []item
	for len(stories) < 30 {
		missing := numStories - len(stories)
		ratio := 1.25
		r := float64(missing) * ratio
		retrieveNum := int(math.Max(r, 1.0))
		retrieveIDs := ids[0:retrieveNum]
		gotStories, err := getPartialStories(retrieveIDs)
		ids = ids[retrieveNum:]
		if err != nil {
			continue
		}
		for _, s := range gotStories {
			stories = append(stories, s)
			if len(stories) == 30 {
				break
			}
		}
	}
	return stories, nil
}

// getPartialStories will spawn number of workers = len(partialIDs)
// return items which are only valid (no error, only story)
func getPartialStories(partialIDs []int) ([]item, error) {
	workerNum := len(partialIDs)
	for seq, id := range partialIDs {

		job := Job{HnID: id, Seq: seq}
		JobQueue <- job
	}
	var tmp []Result
	var partialStories []item

	for i := 0; i < workerNum; i++ {
		r := <-ResultQueue
		if r.Err != nil {
			continue
		}
		if isStoryLink(r.Item) {
			tmp = append(tmp, r)
		}
	}
	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i].Job.Seq < tmp[j].Job.Seq
	})
	for _, res := range tmp {
		partialStories = append(partialStories, res.Item)
	}
	return partialStories, nil
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
