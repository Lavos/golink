package golink

import (
	"log"
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type PhishTank struct {

}

type PhishTankResponse struct {
	Url string
	InDatabase bool
	PhishID int
	PhishDetailPage string
	Verified rune
	VerifiedAt string
	Valid rune
}

// implements Validator
func (s *PhishTank) Validate(checkurl string, response chan Validation) {
	v := url.Values{}
	v.Set("app_key", "517b045104e62d605ff60a33768ed59d92d522fbbcaf643602d56ed9440585be")
	v.Set("url", checkurl)
	v.Set("format", "json")

	resp, err := http.PostForm("http://checkurl.phishtank.com/checkurl/", v)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var phishResponse PhishTankResponse
	err = json.Unmarshal(body, &phishResponse)

	if err != nil {
		log.Printf("count not unmarshal json, %#v", err)
	}

	response <- Validation{
		ValidatorKey: "phishtank",
		URL: checkurl,
		Success: !phishResponse.InDatabase,
	}
}

func init () {
	Validators = append(Validators, &PhishTank{})
}
