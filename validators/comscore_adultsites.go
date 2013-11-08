package validators

import (
	"log"
	"github.com/cznic/kv"
)

type ComscoreAdultSites struct {
	DBLocation string
}

func (c *ComscoreAdultSites) Validate(checkurl string, response chan Validation) {
	db, err := kv.Open(c.DBLocation, &kv.Options{})

	if err != nil {
		log.Fatal("db error: %#v", err)
	}

	defer db.Close()

	v := Validation{
		ValidatorKey: "comscore_adultsites",
		URL: checkurl,
	}

	root, err := GetRootDomain(checkurl)

	if err != nil {
		log.Printf("rootdomain err: %#v", err)
		v.Success = false
	} else {
		var b []byte
		sub, _ := db.Get(b, []byte(root))

		log.Printf("com sub: %#v", sub)

		if sub == nil {
			v.Success = true
		}
	}

	response <- v
}
