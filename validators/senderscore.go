package validators

import (
	"log"
	"net"
	"net/url"
	"net/http"
	"io/ioutil"
	"bytes"
)

type SenderScore struct {

}

func (s *SenderScore) Validate(checkurl string, response chan Validation) {
	u, _ := url.Parse(checkurl)

	var valid = true
	ips, _ := net.LookupIP(u.Host)

	for _, ip := range ips {
		if !s.CheckBlacklist(ip) {
			valid = false;
			break;
		}
	}

	response <- Validation{
		ValidatorKey: "senderscore",
		URL: checkurl,
		Success: valid,
	}
}

func (s *SenderScore) CheckBlacklist(ip net.IP) bool {
	resp, err := http.Get("https://www.senderscore.org/blacklistlookup/index.php?lookup=" + ip.String())

	if err != nil {
		log.Printf("senderscore http error: %#v", err);
		return false
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return bytes.HasPrefix(body, []byte("\n\nnoMessage="))
}
