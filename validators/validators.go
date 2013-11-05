package validators

import (
	"regexp"
	"net/url"
	"errors"
)

var (
	rootDomainRegex = regexp.MustCompile(`((?:[\w\d]*)\.[a-z]{2,}(?:\.[a-z]{2})?)$`)
)

type Configuration struct {
	string
}

type Validation struct {
	ValidatorKey string
	URL string
	Success bool
}

type Validator interface {
	Validate(url string, response chan Validation)
}

type ValidationSet []Validator

// TODO: should probably use PublicSuffix instead
func GetRootDomain(checkurl string) (string, error) {
	u, err := url.Parse(checkurl)

	if err != nil {
		return "", err
	}

	matches := RootDomainRegExp.FindAllString(u.Host, -1)

	if matches == nil {
		return "", errors.New("could not parse domain name against regexp")
	}

	return matches[0], nil
}
