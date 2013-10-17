package validators

import (
	"net/url"
	"net/http"
)

type GoogleSafeBrowsingv1 struct {
	Apikey string
}

func (s *GoogleSafeBrowsingv1) Validate(checkurl string, response chan Validation) {
	v := url.Values{}
	v.Set("apikey", s.Apikey)
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

}
