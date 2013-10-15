package validators

import (
	"time"
)

type SimpleTimer struct {

}

func (s *SimpleTimer) Validate(url string, response chan Validation) {
	time.Sleep(3 * time.Second)
	response <- Validation{
		ValidatorKey: "simple_timer",
		URL: url,
		Success: true,
	}
}

func init () {
	Validators = append(Validators, &SimpleTimer{})
}
