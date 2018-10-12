package check

import (
	"gopkg.in/resty.v1"
)

var client = resty.New().SetRetryCount(3)

func Client() *resty.Client {
	return client
}
