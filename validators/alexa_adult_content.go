package validators

import (
	"log"
	"fmt"
	"time"
	"net/url"
	"net/http"
	"crypto/sha256"
	"crypto/hmac"
	"encoding/base64"
	"encoding/xml"
	"io/ioutil"
)


type UrlInfoResponse struct {
	Response Response `xml:"Response"`
}

type Response struct {
	UrlInfoResult UrlInfoResult `xml:"UrlInfoResult"`
}

type UrlInfoResult struct {
	Alexa AlexaResult `xml:"Alexa"`
}

type AlexaResult struct {
	ContentData ContentData `xml:"ContentData"`
}

type ContentData struct {
	DataUrl string `xml:"DataUrl"`
	AdultContent string `xml:"AdultContent"`
}

type AlexaAdultContent struct {
	AccessKeyID, SecretKey string
}

func (s *AlexaAdultContent) Validate(checkurl string, response chan Validation) {
	va := Validation{
		ValidatorKey: "alexa_adult_content",
		URL: checkurl,
	}

	v := url.Values{}
	v.Set("AWSAccessKeyId", s.AccessKeyID)
	v.Set("Action", "UrlInfo")
	v.Set("ResponseGroup", "AdultContent")
	v.Set("SignatureMethod", "HmacSHA256")
	v.Set("SignatureVersion", "2")
	v.Set("Timestamp", time.Now().UTC().Format("2006-01-02T15:04:05Z"))

	u, err := url.Parse(checkurl)

	if err != nil {
		log.Printf("alexa adult content: %#v", err)
		response <- va
		return
	}

	v.Set("Url", u.Host)

	sig_template := fmt.Sprintf("GET\nawis.amazonaws.com\n/\nAWSAccessKeyId=%s&Action=UrlInfo&ResponseGroup=%s&SignatureMethod=HmacSHA256&SignatureVersion=2&Timestamp=%s&Url=%s", v.Get("AWSAccessKeyId"), v.Get("ResponseGroup"), url.QueryEscape(v.Get("Timestamp")), url.QueryEscape(v.Get("Url")))

	log.Print(sig_template)

	key := []byte(s.SecretKey)
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(sig_template))
	expectedMAC := mac.Sum(nil)

	str := base64.StdEncoding.EncodeToString(expectedMAC)


	v.Set("Signature", str)
	log.Printf("v: %#v", v)

	resp, _ := http.Get("http://awis.amazonaws.com/?" + v.Encode())

	log.Printf("v.encode: %#v", v.Encode())

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	log.Printf("resp: %#v", resp)
	log.Printf("body: %#v", string(body))

	var xmldata UrlInfoResponse
	xml.Unmarshal(body, &xmldata)

	log.Printf("%#v", xmldata)

	va.Success = xmldata.Response.UrlInfoResult.Alexa.ContentData.AdultContent != "yes"
	response <- va
}
