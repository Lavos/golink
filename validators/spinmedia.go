package validators

import (
	"net/url"
	"regexp"
)

var (
	RootDomainRegExp = regexp.MustCompile(`([\w\d]+\.(com|net))$`)
)

type SpinMedia struct {

}

func (s *SpinMedia) Validate(checkurl string, response chan Validation) {
	oo_sites := map[string]bool{
		"spin.com": true,
		"vibe.com": true,
		"stereogum.com": true,
		"idolator.com": true,
		"buzznet.com": true,
		"purevolume.com": true,
		"brooklynvegan.com": true,
		"hypem.com": true,
		"gorillavsbear.net": true,
		"absolutepunk.net": true,
		"concreteloop.com": true,
		"xlr8r.com": true,
		"alterthepress.com": true,
		"directlyrics.com": true,
		"indieshuffle.com": true,
		"popmatters.com": true,
		"prettymuchamazing.com": true,
		"propertyofzack.com": true,
		"underthegunreview.net": true,
		"celebuzz.com": true,
		"egotastic.com": true,
		"thefrisky.com": true,
		"thesuperficial.com": true,
		"gofugyourself.com": true,
		"socialitelife.com": true,
		"pinkisthenewblog.com": true,
		"wwtdd.com": true,
		"fanpop.com": true,
		"crunktastical.net": true,
		"justjared.com": true,
		"justjaredjr.com": true,
		"videogum.com": true,
		"heartsandfoxes.com": true,
	}

	u, _ := url.Parse(checkurl)
	matches := RootDomainRegExp.FindAllString(u.Host, -1)
	valid := false

	if matches != nil && oo_sites[matches[0]] {
		valid = true
	}

	response <- Validation{
		ValidatorKey: "spinmedia",
		URL: checkurl,
		Success: valid,
	}
}
