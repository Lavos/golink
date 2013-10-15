package golink

import (
	"net/http"
	"log"
	"encoding/json"
	// "fmt"
)

type Application struct {
	Port int

	cache *Cache
}

func (a *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result := a.cache.Hit(r.URL.String()[1:])

	log.Printf("result: %#v", result)

	b, _ := json.Marshal(result)
	w.Write(b)
}

func (a *Application) Run() {
	log.Print("started server")
	http.ListenAndServe(":8000", a)
}

func NewApplication (port int) *Application {
	a := &Application{
		Port: port,
		cache: newCache(),
	}

	return a
}
