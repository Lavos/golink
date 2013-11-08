package golink

import (
	"log"
	"encoding/json"
	"github.com/Lavos/golink/validators"
)

type Cache struct {
	links map[string]*Result

	blacklist, whitelist *validators.ValidationSet
	fillBlack, fillWhite chan validators.Validation

	query chan cacheRequest
}

type cacheRequest struct {
	URL string
	response chan []byte
}

type Result struct {
	URL string `json:"url,omitempty"`
	Determination string `json:"determination,omitempty"`
	Reason string `json:"reason,omitempty"`
	Status string `json:"status"`
	WhiteList map[string]bool `json:"whitelist"`
	BlackList map[string]bool `json:"blacklist"`
}


func (c *Cache) Hit(url string) []byte {
	req := cacheRequest{
		URL: url,
		response: make(chan []byte),
	}

	c.query <- req
	return <-req.response
}

func (c *Cache) validateAll(url string) {
	for _, v := range *c.blacklist {
		c.validateAgainst(v, url, c.fillBlack)
	}

	for _, v := range *c.whitelist {
		c.validateAgainst(v, url, c.fillWhite)
	}
}

func (c *Cache) validateAgainst(v validators.Validator, url string, responsechan chan validators.Validation) {
	go v.Validate(url, responsechan)
}

func (c *Cache) checkComplete(r *Result) {
	if len(r.BlackList) + len(r.WhiteList) != len(*c.blacklist) + len(*c.whitelist) {
		log.Print("not done.")
		return
	}

	log.Print("DONE!")
	r.Status = "done"

	// blacklist logic
	for _, b := range r.BlackList {
		if !b {
			r.Determination = "remove"
			r.Reason = "blacklist failure"
			return
		}
	}

	// whitelist logic
	for _, w := range r.WhiteList {
		if w {
			r.Determination = "unmodify"
			r.Reason = "whitelist success"
			return
		}
	}

	r.Determination = "nofollow"
}

func (c *Cache) run() {
	c.links = make(map[string]*Result)

	for {
		select {
		case request := <-c.query:
			if _, ok := c.links[request.URL]; !ok {
				c.links[request.URL] = &Result{
					URL: request.URL,
					Status: "validating",
					WhiteList: make(map[string]bool),
					BlackList: make(map[string]bool),
				}

				c.validateAll(request.URL)
			}

			result := c.links[request.URL]
			b, _ := json.Marshal(result)
			request.response <- b

		case v := <-c.fillBlack:
			log.Print("FILLBLACK")

			r := c.links[v.URL]
			r.BlackList[v.ValidatorKey] = v.Success;
			c.checkComplete(r)

		case v := <-c.fillWhite:
			log.Print("FILLWHITE")

			r := c.links[v.URL]
			r.WhiteList[v.ValidatorKey] = v.Success;
			c.checkComplete(r)
		}
	}
}

func newCache (blacklist, whitelist *validators.ValidationSet) *Cache {
	c := &Cache{
		query: make(chan cacheRequest),
		fillWhite: make(chan validators.Validation),
		fillBlack: make(chan validators.Validation),
		blacklist: blacklist,
		whitelist: whitelist,
	}

	go c.run()

	return c
}
