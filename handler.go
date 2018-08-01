package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/sourcegraph/checkup"
)

type Notify struct {
}

func homeHandler(c echo.Context) error {
	//may do redirect later?
	return c.String(http.StatusOK, "home page")
}

// if provided name only, do a active check
// otherwise do a info register and active check
func checkHandler(c echo.Context) error {
	x := &checkup.HTTPChecker{}
	if err := c.Bind(x); err != nil {
		log.Println("checkhandler", err)
		return err
	}
	return c.JSON(http.StatusCreated, x)
}

func notifyHandler(c echo.Context) error {
	log.Println("starting notify")
	n := &checkup.Qianbao{}
	if err := c.Bind(n); err != nil {
		log.Println("notify err", err)
		return err
	}
	return c.JSON(http.StatusCreated, n)
}
