package main

import (
	"fmt"
	"net/http"
	"strings"
	"wen/hook-api/upstream"

	"github.com/chinglinwen/log"

	"github.com/chinglinwen/checkup"
	"github.com/labstack/echo"

	"wen/svc-d/fetch"
)

func homeHandler(c echo.Context) error {
	//may do redirect later?
	return c.String(http.StatusOK, "home page")
}

// if provided name only, do a active check
// otherwise do a info register and active check
//
// if provide ip, use that ip
// fetch config from project center, if non-provided, use default?
func checkHandler(c echo.Context) error {
	appname := c.FormValue("appname")
	env := c.FormValue("env")

	/* 	x := &fetch.ProjectCheck{}
	   	if err := c.Bind(x); err != nil {
	   		log.Println("checkhandler bind", err)
	   		return err
	   	} */

	p, err := fetch.Fetch(env, appname)
	if err != nil {
		log.Println("fetch project err: ", err)
		return err
	}
	err = p.Check()
	if err != nil {
		log.Printf("%v check err: %v\n", p.Name, err)
		return err
	}
	out := fmt.Sprintf("%v check ok", p.Name)
	log.Println(out)
	return c.JSON(http.StatusOK, out)
}

func notifyHandler(c echo.Context) error {
	r := &checkup.Result{}
	if err := c.Bind(r); err != nil {
		log.Println("notify handler bind err", err)
		return err
	}
	log.Println("notify: ", r)

	name, ns := getNamespace(r.Title)
	ok, err := upstream.ChangeState(r.Endpoint, name, ns, "0")
	if err != nil {
		log.Printf("notify:  %v %v, change upstream state, err: %v", r.Title, r.Endpoint, err)
		return err
	}
	if ok == false {
		log.Printf("notify:  %v %v, change upstream state failure", r.Title, r.Endpoint)
	}
	return c.String(http.StatusOK, "notify page")
}

func getNamespace(title string) (name, ns string) {
	if title == "" {
		return
	}
	s := strings.Split(title, "#")
	name = s[0]
	if len(s) >= 2 {
		ns = s[1]
	}
	return
}
