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
		e := fmt.Sprintf("fetch project err: ", err)
		return c.JSONPretty(400, E(1, e, "error"), " ")
	}
	err = p.Check()
	if err != nil {
		e := fmt.Sprintf("api check for %v err: %v", p.Name, err)
		return c.JSONPretty(400, E(2, e, "error"), " ")
	}
	out := fmt.Sprintf("api check for %v ok", p.Name)
	log.Println(out)

	return c.JSONPretty(http.StatusOK, E(0, out, "ok"), " ")
}

func E(code int, msg, status string) map[string]interface{} {
	log.Println(msg)
	return map[string]interface{}{
		"code":   code,
		"msg":    msg,
		"status": status,
	}
}

func notifyHandler(c echo.Context) error {
	r := &checkup.Result{}
	if err := c.Bind(r); err != nil {
		e := fmt.Sprintf("notify handler bind err", err)
		return c.JSONPretty(400, E(1, e, "error"), " ")
	}
	log.Println("notify: ", r)

	name, ns := getNamespace(r.Title)
	ok, err := upstream.ChangeState(r.Endpoint, name, ns, "0")
	if err != nil {
		e := fmt.Sprintf("notify:  %v %v, change upstream state, err: %v", r.Title, r.Endpoint, err)
		return c.JSONPretty(400, E(2, e, "error"), " ")
	}

	if ok == false {
		e := fmt.Sprintf("notify:  %v %v, change upstream state failure", r.Title, r.Endpoint)
		return c.JSONPretty(400, E(3, e, "error"), " ")
	}

	return c.JSONPretty(http.StatusOK, E(0, "notify ok", "ok"), " ")
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
