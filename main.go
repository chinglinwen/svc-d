package main

import (
	"fmt"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/sourcegraph/checkup"
)

func main() {

	c := checkup.Checkup{
		Checkers: []checkup.Checker{
			checkup.HTTPChecker{
				Name:     "Website",
				URL:      "http://www.baidu.com",
				Attempts: 5,
			},
		},
		Storage: checkup.FS{
			Dir:         "./data",
			CheckExpiry: 7 * 24 * time.Hour,
		},
		Notifier: Wechat{},
	}
	spew.Dump(c)
	/*
		// perform a checkup
		results, err := c.CheckAndStore()
		if err != nil {
			log.Fatal(err)
		}
		for _, result := range results {
			fmt.Println(result)
		}
	*/
	fmt.Println("start")
	err := c.CheckAndStore()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("done")
}
