package golink

import (
	"log"
	"github.com/Lavos/golink/validators"
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

	blacklist, whitelist *validators.ValidationSet
	query chan cacheRequest
	fill chan validators.Validation
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
	for _, v := range *c.blacklist {
		c.validateAgainst(v, url)
	}

	for _, v := range *c.whitelist {
		c.validateAgainst(v, url)
	}
}

func (c *Cache) validateAgainst(v validators.Validator, url string) {
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

				c.validateAll(request.URL)
			}

			request.response <- c.urls[request.URL]

		case v := <-c.fill:
			log.Printf("got from validator: %#v", v)
			current := c.urls[v.URL]

			current.Validations[v.ValidatorKey] = v.Success

			if len(current.Validations) == len(*c.blacklist) {
				current.Status = "done"

				d := "unmodify"

				for _, b := range current.Validations {
					if !b {
						d = "remove"
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

func newCache (blacklist, whitelist *validators.ValidationSet) *Cache {
	c := &Cache{
		urls: make(map[string]Result),
		query: make(chan cacheRequest),
		fill: make(chan validators.Validation),
		blacklist: blacklist,
		whitelist: whitelist,
	}

	go c.run()

	return c
}
