package golink

import (
	"net/http"
	"log"
	"encoding/json"
	// "fmt"
)

type Application struct {
	Port int

	cache map[string]Result
	hitcache chan cacheRequest
	results chan Result
	queue chan string
}

type cacheRequest struct {
	URL string
	responseChan chan cacheResponse
}

type cacheResponse struct {
	success bool
	Result Result
}

func (a *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := cacheRequest{r.URL.String(), make(chan cacheResponse)}
	a.hitcache <- req

	pair := <-req.responseChan

	if pair.success {
		log.Print("Hey, I have this value!")

		b, _ := json.Marshal(pair.Result)
		w.Write(b)
	} else {
		log.Print("wow, this is brand new.")

		b, _ := json.Marshal(Result{
			URL: req.URL,
			Status: "queued",
		})
		w.Write(b)

		a.queue <- r.URL.String()
	}
}

func (a *Application) Run() {
	log.Print("started server")
	go http.ListenAndServe(":8000", a)

	for {
		select {
		case result := <-a.results:
			a.cache[result.URL] = result
		case h := <-a.hitcache:
			log.Printf("%#v", a.cache)

			if r, ok := a.cache[h.URL]; ok {
				log.Printf("%#v", ok)

				h.responseChan <- cacheResponse{
					success: true,
					Result: r,
				}
			} else {
				h.responseChan <- cacheResponse{
					success: false,
				}
			}
		}
	}
}

func NewApplication (port int) *Application {
	a := &Application{
		Port: port,
		cache: make(map[string]Result),
		hitcache: make(chan cacheRequest),
		results: make(chan Result),
		queue: make(chan string, 100),
	}

	newQueue(a.queue, a.results)

	return a
}
