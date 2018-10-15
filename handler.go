package main

import (
	"fmt"
	"net/http"
	"strings"
	"wen/hook-api/upstream"

	"github.com/chinglinwen/log"

	"github.com/chinglinwen/checkup"
	"github.com/labstack/echo"

	"wen/svc-d/check"
	"wen/svc-d/notice"
)

/* func homeHandler(c echo.Context) error {
	//may do redirect later?
	return c.String(http.StatusOK, "home page")
} */

// if provided name only, do a active check
// otherwise do a info register and active check
//
// if provide ip, use that ip
// fetch config from project center, if non-provided, use default?
func checkHandler(c echo.Context) error {
	appname := c.FormValue("appname")
	env := c.FormValue("env")
	if appname == "" {
		return c.JSONPretty(400, E(1, "appname not provided", "error"), " ")
	}

	p, err := check.Fetch(env, appname)
	if err != nil {
		e := fmt.Sprintf("fetch project err: %v", err)
		return c.JSONPretty(400, E(1, e, "error"), " ")
	}
	r, err := p.Check()
	if err != nil {
		e := fmt.Sprintf("fetch project err: %v", err)
		return c.JSONPretty(400, E(2, e, "error"), " ")
	}
	log.Debug.Println("result", r)

	var errfound bool

	for _, v := range r {
		if !v.Healthy {
			errfound = true
			log.Printf("%v, healthy: %v\n", v.Endpoint, v.Healthy)
		}
	}
	if errfound {
		out := EData(3, "check failure found", "error", Result2Map(r))
		return c.JSONPretty(400, out, " ")
	}

	out := fmt.Sprintf("api check for %v ok", p.Name)
	log.Println(out)

	return c.JSONPretty(http.StatusOK, E(0, out, "ok"), " ")
}

func E(code int, msg, status string) map[string]interface{} {
	log.Println(msg)
	return map[string]interface{}{
		"code":    code,
		"message": msg,
		"status":  status,
	}
}
func Attempt2Error(times checkup.Attempts) (err string) {
	i := 0
	for _, v := range times {
		if i == 0 {
			err += v.Error
		} else {
			err += "\n" + v.Error
		}
		return
	}
	return
}

func Result2Map(r []checkup.Result) (data []map[string]interface{}) {
	for _, v := range r {
		if !v.Healthy {
			name, _ := getNamespace(v.Title)
			ip, port := check.Endpoint2ip(v.Endpoint)
			d := map[string]interface{}{
				"name":   name,
				"ip":     ip,
				"port":   port,
				"status": v.Status(),
				"reason": Attempt2Error(v.Times),
			}
			data = append(data, d)
		}
	}
	return data
}

func EData(code int, msg, status string, data []map[string]interface{}) map[string]interface{} {
	log.Println(msg)
	return map[string]interface{}{
		"code":    code,
		"message": msg,
		"status":  status,
		"data":    data,
	}
}

func notifyHandler(c echo.Context) error {
	r := &checkup.Result{}
	if err := c.Bind(r); err != nil {
		e := fmt.Sprintf("notify handler bind err %v", err)
		return c.JSONPretty(400, E(1, e, "error"), " ")
	}
	log.Println("notify: ", r)

	name, ns := getNamespace(r.Title)
	endpoint := r.Endpoint

	if *testproject != "" {
		if *testproject != name {
			e := fmt.Sprintf("notify:  %v %v, it's not test project, skip change", name, endpoint)
			return c.JSONPretty(http.StatusOK, E(0, e, "error"), " ")
		}
	}

	// get config from project center.
	config, err := check.FetchConfig(name)
	if err != nil {
		log.Println("fetch config for notify err", err)
	}

	content := fmt.Sprintf("%v %v, changed upstream state", name, endpoint)

	// do send alert here, if receiver is not empty
	// things come to here is not ok.

	// no need for now
	if config.AlertReceiver != "" {
		reply, err := notice.Send(config.AlertReceiver, content, "down", "10m")
		if err != nil {
			log.Printf("send alert err: %v, resp: %v\n", err, reply)
		} else {
			log.Printf("send alert ok, resp: %v\n", reply)
		}
	}

	// extra alert if alertall is setted.
	// will do alert in checkup instead
	/* 	if *alertAll {
		reply, err := notice.Send(*defaultReceiver, content)
		if err != nil {
			log.Printf("send alertall err: %v, resp: %v\n", err, reply)
		} else {
			log.Printf("send alertall ok, resp: %v\n", reply)
		}
	} */

	if config.AutoDisable != "on" {
		msg := "no setting of auto change state for upstream for " + name + " ns " + ns
		return c.JSONPretty(http.StatusOK, E(0, msg, "ok"), " ")
	}

	ok, err := upstream.ChangeState(r.Endpoint, name, ns, "0")
	if err != nil {
		// do send alert here, things err
		content := fmt.Sprintf("%v %v, change upstream state err: %v", name, endpoint, err)
		reply, err := notice.Send(*defaultReceiver, content, "error", "")
		if err != nil {
			log.Printf("send alert err: %v, resp: %v\n", err, reply)
		} else {
			log.Printf("send alert ok, resp: %v\n", reply)
		}

		e := fmt.Sprintf("notify:  %v %v, change upstream state, err: %v", name, endpoint, err)
		return c.JSONPretty(400, E(2, e, "error"), " ")
	}

	if ok == false {
		e := fmt.Sprintf("notify:  %v %v, change upstream state failure", name, endpoint)
		return c.JSONPretty(400, E(3, e, "error"), " ")
	}

	msg := fmt.Sprintf("notify:  %v %v, change upstream state ok", name, endpoint)
	return c.JSONPretty(http.StatusOK, E(0, msg, "ok"), " ")
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
