package validators

import (
	"log"
	"github.com/cznic/kv"
)

type AlexaTopSites struct {
	DBLocation string
}

func (a *AlexaTopSites) Validate(checkurl string, response chan Validation) {
	db, err := kv.Open(a.DBLocation, &kv.Options{})

	if err != nil {
		log.Fatal("db error: %#v", err)
	}

	defer db.Close()

	v := Validation{
		ValidatorKey: "alexa_topsites",
		URL: checkurl,
	}

	root, err := GetRootDomain(checkurl)

	if err != nil {
		v.Success = false
	} else {
		var b []byte
		sub, _ := db.Get(b, []byte(root))

		if sub != nil {
			v.Success = true
		}
	}

	response <- v
}
