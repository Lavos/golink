package main

import (
	"log"
	"os"
	"encoding/csv"
	"github.com/cznic/kv"
)

var (
	domain string
)

func main() {
	db, err := kv.Create("alexa_topsites.kv", &kv.Options{})

	if err != nil {
		log.Fatalf("%#v", err)
	}

	csv := csv.NewReader(os.Stdin)

	for {
		record, err := csv.Read()

		if err != nil {
			break
		}

		log.Printf("adding: %#v", record[1])

		db.Set([]byte(record[1]), []byte(record[0]))
	}

	db.Close()
}
