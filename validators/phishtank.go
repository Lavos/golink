package validators

import (
	"log"
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type PhishTank struct {
	AppKey string
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

func (s *PhishTank) Validate(checkurl string, response chan Validation) {
	v := url.Values{}
	v.Set("app_key", s.AppKey)
	v.Set("url", checkurl)
	v.Set("format", "json")

	resp, err := http.PostForm("http://checkurl.phishtank.com/checkurl/", v)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var phishResponse PhishTankResponse
	err = json.Unmarshal(body, &phishResponse)

	if err != nil {
		log.Printf("could not unmarshal json, %#v", err)
	}

	response <- Validation{
		ValidatorKey: "phishtank",
		URL: checkurl,
		Success: !phishResponse.InDatabase,
	}
}

func init () {

}
