package validators

import (

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
