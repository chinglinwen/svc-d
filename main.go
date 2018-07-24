package main

import (
	"fmt"
	"log"
	"time"
	"wen/svc-d/config"

	"github.com/davecgh/go-spew/spew"
	"github.com/sourcegraph/checkup"
)

//define a global variable
//add new check, update it, and store the config as file(update config)

func main() {
	log.Println("starting...")

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
	cc := config.New("test.json")
	cc.Checkup = c

	log.Println("start save")
	err := cc.Save()
	if err != nil {
		log.Println("save error")
		log.Fatal(err)
	}
	spew.Dump(cc)
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
	err = c.CheckAndStore()
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("done")

}
