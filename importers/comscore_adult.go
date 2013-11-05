package main

import (
	"log"
	"fmt"
	"strings"
	"github.com/cznic/kv"
)

func main () {
	db, err := kv.Create("comscore_adult.kv", &kv.Options{})
	defer db.Close()

	if err != nil {
		log.Fatalf("Could not create database: %#v")
	}

	var domain string

	for {
		_, err := fmt.Scanln(&domain)

		if err != nil {
			break
		}

		log.Printf("adding: %#v", domain)
		db.Set([]byte(strings.ToLower(domain)), []byte("1"))
	}
}
