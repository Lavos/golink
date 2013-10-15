package golink

import (
	"log"
)

type Result struct {
	URL string `json:"url"`
	Determination string `json:"determination,omitempty"`
	Reason string `json:"reason,omitempty"`
	Status string `json:"status"`
	Validations map[string]bool `json:"validations"`
}

type Cache struct {
	urls map[string]Result

	query chan cacheRequest
	fill chan Validation
}

type cacheRequest struct {
	URL string
	response chan Result
}


func (c *Cache) Hit(url string) Result {
	req := cacheRequest{
		URL: url,
		response: make(chan Result),
	}

	c.query <- req
	result := <-req.response
	return result
}

func (c *Cache) validateAll(url string) {
	for _, v := range Validators {
		c.validateAgainst(v, url)
	}
}

func (c *Cache) validateAgainst(v Validator, url string) {
	go v.Validate(url, c.fill)
}

func (c *Cache) run() {
	for {
		select {
		case request := <-c.query:
			if _, ok := c.urls[request.URL]; !ok {
				c.urls[request.URL] = Result{
					URL: request.URL,
					Validations: make(map[string]bool),
					Status: "validating",
				}

			}

			request.response <- c.urls[request.URL]

		case v := <-c.fill:
			log.Printf("got from validator: %#v", v)
			current := c.urls[v.URL]

			current.Validations[v.ValidatorKey] = v.Success

			if len(current.Validations) == len(Validators) {
				current.Status = "done"

				d := "accepted"

				for _, b := range current.Validations {
					if !b {
						d = "rejected"
						break
					}
				}

				log.Printf("%#v", d)
				current.Determination = d
			}

			c.urls[v.URL] = current
		}
	}
}

func newCache () *Cache {
	c := &Cache{
		urls: make(map[string]Result),
		query: make(chan cacheRequest),
		fill: make(chan Validation),
	}

	go c.run()

	return c
}
