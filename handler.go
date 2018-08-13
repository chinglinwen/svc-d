package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/chinglinwen/log"
	"github.com/go-resty/resty"

	"github.com/chinglinwen/checkup"
	"github.com/labstack/echo"
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
	r := &checkup.Result{}
	if err := c.Bind(r); err != nil {
		log.Println("notify handler bind err", err)
		return err
	}

	log.Println("notify: ", r)
	ok, err := ChangeState(r.Endpoint, r.Title)
	if err != nil {
		log.Printf("notify:  %v %v, change upstream state, err: %v", r.Title, r.Endpoint, err)
		return err
	}
	if ok == false {
		log.Printf("notify:  %v %v, change upstream state failure", r.Title, r.Endpoint)
	}
	return c.String(http.StatusOK, "notify page")
}

// ChangeState change project specific ip state, remove item from nginx
// The logic may need to distinguish VM and docker
// We currently check based on docker first.
//
// Upstream will make it disabled ( need rethink?)
// Service recovery need human manual operation.
func ChangeState(endpoint, title string) (bool, error) {
	ip, port := endpoint2ip(endpoint)
	resp, err := resty.R().
		SetFormData(map[string]string{
			"appname": title,
			"ip":      ip,
			"port":    port,
			"state":   "0", // int 1:up or 0:down
		}).
		Post(*upstreamnChangeAPI)

	if err != nil {
		return false, err
	}
	state, err := parseState(resp.Body())
	if state != true {
		log.Println("ChangeState resp: ", string(resp.Body()))
	}
	return state, err
}

func parseState(body []byte) (state bool, err error) {
	var result []interface{}
	err = json.Unmarshal(body, &result)
	if err != nil || len(result) == 0 {
		return
	}
	if state, _ = result[0].(bool); state != true {
		return
	}
	return true, err
}

func endpoint2ip(e string) (ip, port string) {
	str := strings.Split(e, "/")
	var s string
	if len(str) > 2 {
		s = str[2]
	} else if len(str) == 1 {
		s = str[0]
	}
	ipport := strings.Split(s, ":")
	if len(ipport) == 2 {
		ip, port = ipport[0], ipport[1]
	}
	return
}
