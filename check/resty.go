package check

import (
	"gopkg.in/resty.v1"
)

var DEBUG bool
var client = resty.New().SetRetryCount(3)

func Client() *resty.Client {
	if DEBUG {
		return client.SetDebug(true)
	}
	return client
}
