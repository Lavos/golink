package golink

import (
	"log"
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type GoogleSafeBrowsingv1 struct {

}

func (s *GoogleSafeBrowsingv1) Validate(checkurl string, response chan Validation) {
	v := url.Values{}
	v.Set("apikey", "ABQIAAAANzZNsoU24Ofp83gTQlx9xhTdfLkavSFRlhqu8kud5wNuQRJudg")
	v.Set("client", "golink")
	v.Set("appver", "0.0.0.1")
	v.Set("pver", "3.0")
	v.Set("url", checkurl)

	resp, _ := http.Get("https://sb-ssl.google.com/safebrowsing/api/lookup?" + v.Encode())

	response <- Validation{
		ValidatorKey: "google",
		URL: checkurl,
		Success: resp.StatusCode == 200 || resp.StatusCode == 204,
	}
}

func init () {
	Validators = append(Validators, &GoogleSafeBrowsingv1{})
}
