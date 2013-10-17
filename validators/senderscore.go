package validators

import (
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
	resp, _ := http.Get("https://www.senderscore.org/blacklistlookup/index.php?lookup=" + ip.String())
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return bytes.HasPrefix(body, []byte("\n\nnoMessage="))
}

func init () {

}
