package main

import (
	"fmt"

	"wen/svc-d/config"

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
	go start()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", homeHandler)
	e.GET("/check", checkHandler)
	e.GET("/notify", notifyHandler)

	err := e.Start(":1323")
	fmt.Println("fatal", err)
	//e.Logger.Fatal()

	fmt.Println("exit")

}
