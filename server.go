package golink

import (
	"net/http"
	"log"
	"os"
	"github.com/Lavos/golink/validators"
)

type Application struct {
	Address string

	cache *Cache
}

type Server struct {
	cache *Cache
}

type Configuration struct {
	Address, GoogleAPIKey, AlexaAccessKeyID, AlexaSecretKey, PhishTankAppKey, AlexaTopSitesDB, ComscoreAdultSitesDB string
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	if q != "" {
		result := s.cache.Hit(q)
		w.Write(result)
	} else {
		w.WriteHeader(400)
	}
}

func (a *Application) Run() {
	log.Print("started server")
	s := &Server{ a.cache }

	go http.ListenAndServe(a.Address, s)

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
