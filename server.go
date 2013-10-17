package golink

import (
	"net/http"
	"log"
	"encoding/json"
	"github.com/Lavos/golink/validators"
)

type Application struct {
	Address string

	cache *Cache
}

type Configuration struct {
	Address, GoogleAPIKey, AlexaAccessKeyID, AlexaSecretKey, PhishTankAppKey string
}

func (a *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result := a.cache.Hit(r.URL.String()[1:])

	log.Printf("result: %#v", result)

	b, _ := json.Marshal(result)
	w.Write(b)
}

func (a *Application) Run() {
	log.Print("started server")
	http.ListenAndServe(a.Address, a)
}

func NewApplication (c *Configuration) *Application {
	v := &validators.ValidationSet{
		&validators.GoogleSafeBrowsingv1{ c.GoogleAPIKey },
		&validators.Alexa{ AccessKeyID: c.AlexaAccessKeyID, SecretKey: c.AlexaSecretKey },
		&validators.PhishTank{ c.PhishTankAppKey },
		&validators.SenderScore{ },
	}

	a := &Application{
		Address: c.Address,
		cache: newCache(v),
	}

	return a
}
