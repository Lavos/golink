package validators

import (

)

var (
	Validators []Validator
)

type Validation struct {
	ValidatorKey string
	URL string
	Success bool
}

type Validator interface {
	Validate(url string, response chan Validation)
}
