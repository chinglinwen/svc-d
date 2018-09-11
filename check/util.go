package check

import (
	"time"

	"github.com/chinglinwen/checkup"
)

func GetCheckWithConfig(name, ip, port string, config ProjectCheck) (check checkup.Checker) {
	if config.Type == "http" {
		c := checkup.HTTPChecker{
			Name: name,
			URL:  "http://" + ip + ":" + port + config.URI,
		}
		if config.StatusCode == 0 {
			c.UpStatus = 452 //above 500 consider error
		} else {
			c.UpStatus = config.StatusCode
		}
		if config.Attempts != 0 {
			c.Attempts = config.Attempts
		}
		if config.AttemptSpacing != 0 {
			c.AttemptSpacing = time.Duration(config.AttemptSpacing) * time.Second
		}
		if config.MustContain != "" {
			c.MustContain = config.MustContain
		}
		check = c
	} else {
		c := checkup.TCPChecker{
			Name: name, // notify api will use it
			URL:  ip + ":" + port,
		}
		if config.Attempts != 0 {
			c.Attempts = config.Attempts
		}
		if config.AttemptSpacing != 0 {
			c.AttemptSpacing = time.Duration(config.AttemptSpacing) * time.Second
		}
		if config.Timeout != 0 {
			c.Timeout = time.Duration(config.Timeout) * time.Second
		}
		check = c
	}
	return
}
