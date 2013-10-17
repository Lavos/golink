package main

import (
	"os"
	"github.com/Lavos/golink"
	"log"
	"encoding/json"
	"flag"
)

var (
	config_filename = flag.String("c", "", "filename of json configuration file")
)

func main () {
	// configuration JSON
	flag.Parse()
	log.Printf("config filename: %#v", *config_filename)

	config_file, err := os.Open(*config_filename)
	defer config_file.Close()

	if err != nil {
		log.Fatal(err)
	}

	data := make([]byte, 2049)
	n, err := config_file.Read(data)

	if err != nil {
		log.Fatal(err)
	}

	var config golink.Configuration
	err = json.Unmarshal(data[:n], &config)

	if err != nil {
		log.Fatal(err)
	}

	app := golink.NewApplication(&config)
	app.Run()
}
