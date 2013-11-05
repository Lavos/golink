package main

import (
	"log"
	"github.com/cznic/kv"
)


func main() {
	db, _ := kv.Open("comscore_adult.kv", &kv.Options{})
	defer db.Close()

	size, _ := db.Size()
	log.Printf("%#v", size)

	buf := make([]byte, 20)
	sub := make([]byte, 0)

	sub, _ = db.Get(buf, []byte("cutepet.org"))
	log.Printf("%#v", string(sub))

	sub, _ = db.Get(buf, []byte("livestrong.com"))
	log.Printf("%#v", string(sub))
}
