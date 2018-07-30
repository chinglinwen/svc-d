package main

import (
	"fmt"
	"log"

	"wen/svc-d/config"

	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var conf *config.Config

func init() {
	conf = config.New("test.json") //try the item, project based ?  why not just name?
}

//define a global variable
//add new check, update it, and store the config as file(update config)

func main() {
	fmt.Println("starting...")
	log.Println("starting... log")

	/*
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
			Notifier: &Wechat{URL: "test url"},
		}
		cc.Checkup = c


		fmt.Println("start save")
		err := cc.Save()
		if err != nil {
			fmt.Println("save error", err)
			os.Exit(1)
		}
	*/

	//spew.Dump(cc)
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

	c := Check{conf}
	/*
		c.Notifier = checkup.Qianbao{Name: "qianbao", Channel: "test url"}

		err := c.Save()
		if err != nil {
			fmt.Println("save error", err)
			os.Exit(1)
		}
	*/
	spew.Dump(c)
	/*
		check := checkup.HTTPChecker{
			Name:     "Website1",
			URL:      "http://www.baidu.com",
			Attempts: 5,
		}

		tcpcheck := checkup.TCPChecker{
			Name:     "tcp1",
			URL:      "220.181.111.188:80",
			Attempts: 5,
		}

		c.Checkup.Checkers = append(c.Checkup.Checkers, check, tcpcheck)
		spew.Dump(c)

		err := c.CheckAndStore()
		if err != nil {
			fmt.Println("checkandstore", err)
			os.Exit(1)
			return
		}
	*/

	fmt.Println("done")

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", homeHandler)
	e.GET("/check", checkHandler)
	e.POST("/notify", notifyHandler)

	err := e.Start(":1323")
	fmt.Println("fatal", err)
	//e.Logger.Fatal()

}
