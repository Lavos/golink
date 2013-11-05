package golink

import (
	"net/http"
	"log"
	"os"
	"encoding/json"
	"github.com/Lavos/golink/validators"
)

type Application struct {
	Address string

	cache *Cache
}

type Configuration struct {
	Address, GoogleAPIKey, AlexaAccessKeyID, AlexaSecretKey, PhishTankAppKey, AlexaTopSitesDB, ComscoreAdultSitesDB string
}

func (a *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result := a.cache.Hit(r.URL.String()[1:])

	log.Printf("result: %#v", result)

	b, _ := json.Marshal(result)
	w.Write(b)
}

func (a *Application) Run() {
	log.Print("started server")
	go http.ListenAndServe(a.Address, a)

	awaitQuitKey()
}

func awaitQuitKey() {
	var buf [1]byte
	for {
		_, err := os.Stdin.Read(buf[:])
		if err != nil || buf[0] == 'q' {
			return
		}
	}
}

func NewApplication (c *Configuration) *Application {
	blacklist := &validators.ValidationSet{
		&validators.GoogleSafeBrowsingv1{ c.GoogleAPIKey },
		&validators.AlexaAdultContent{ AccessKeyID: c.AlexaAccessKeyID, SecretKey: c.AlexaSecretKey },
		&validators.ComscoreAdultSites{ c.ComscoreAdultSitesDB },
		&validators.PhishTank{ c.PhishTankAppKey },
		&validators.SenderScore{ },
	}

	whitelist := &validators.ValidationSet{
		&validators.AlexaTopSites{ c.AlexaTopSitesDB },
		&validators.SpinMedia{ },
	}

	a := &Application{
		Address: c.Address,
		cache: newCache(blacklist, whitelist),
	}

	return a
}
