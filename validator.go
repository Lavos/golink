package golink

import (
	"log"
	"time"
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/Lavos/golink/validators/"
)

var (
	Validators []*Validator
)

type Validation struct {
	ValidatorKey string
	URL string
	Success bool
}

type Validator interface {
	Validate(url string, response chan Validation)
}


type SimpleTimer struct {
	Key string
}

// implements Validator
func (s *SimpleTimer) Validate(url string, response chan Validation) {
	log.Printf("validate: %#v", url)

	time.Sleep(3 * time.Second)
	response <- Validation{
		ValidatorKey: s.Key,
		URL: url,
		Success: true,
	}
}
