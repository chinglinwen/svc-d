package upstream

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chinglinwen/log"
	resty "gopkg.in/resty.v1"
)

var (
	UpstreamBase       string
	UpstreamAddAPI     string
	UpstreamDelAPI     string
	UpstreamnChangeAPI string
)

func Init(UpstreamBase string) {
	UpstreamAddAPI = UpstreamBase + "/wk_add_upstream/"
	UpstreamDelAPI = UpstreamBase + "/wk_deleted_upstream/"
	UpstreamnChangeAPI = UpstreamBase + "/up_nginx_state/"

	//UpstreamAllAPI    = UpstreamBase + "/get_upstream_all_instance/"
	//UpstreamSingleAPI = UpstreamBase + "/get_nginx_all/"
}

type Upstream struct {
	Name      string
	Namespace string
	State     string // "0 disable, 1 enable"
	Env       string // "pre|pro|qa"
	IP        string // single ip only for now
	Port      string
	IsDocker  string // "0 is vm ,1 is docker"
	NginxGrp  string // nginx group, "BJ-SH"

	// "fail_timeout": "int default 30"
	// "weight": "int default 1000"
}

func (u *Upstream) AddValidate() error {
	if u.Name == "" {
		return fmt.Errorf("wk_name not provided")
	}
	if u.Namespace == "" {
		return fmt.Errorf("namespace not provided")
	}
	if u.State == "" {
		return fmt.Errorf("state not provided")
	}
	if u.IP == "" {
		return fmt.Errorf("ip not provided")
	}
	if u.Port == "" {
		return fmt.Errorf("port not provided")
	}
	if u.Env == "" {
		return fmt.Errorf("env not provided")
	}
	if u.IsDocker == "" {
		return fmt.Errorf("isdocker not provided")
	}
	if u.NginxGrp == "" {
		return fmt.Errorf("nginxgrp not provided")
	}
	return nil
}

func (u *Upstream) Add() (err error) {
	if err = u.AddValidate(); err != nil {
		return fmt.Errorf("parameter for add validate err: %v", err)
	}
	return retry("add", u.add, 3)
}

func retry(name string, f func() error, n int) (err error) {
	for i := 1; i <= n; i++ {
		err = f()
		if err != nil {
			log.Printf("%v failed %v times\n", name, i)
			continue
		}
		log.Printf("%v is ok now\n", name)
		return nil
	}
	log.Printf("end of trying %v, after %v times, err\n", name, n, err)
	return
}

func (u *Upstream) add() error {
	resp, err := resty.SetRetryCount(3).
		//SetDebug(true).
		R().SetFormData(map[string]string{
		"wk_name":   u.Name,
		"namespace": u.Namespace,
		"state":     u.State,
		"ip_list":   u.IP,
		"port":      u.Port,
		"env":       u.Env,
		"is_docker": u.IsDocker,
		"nginx":     u.NginxGrp,
	}).
		Post(UpstreamAddAPI)

	if err != nil {
		return err
	}
	log.Println("resp: ", strings.Replace(limit(resp.Body()), "\n", "", -1))

	state, err := parseState(resp.Body())
	if err != nil {
		return fmt.Errorf("upstream add err resp: %v", limit(resp.Body()))
	}
	if state != true {
		return fmt.Errorf("upstream add result failed resp: %v", limit(resp.Body()))
	}
	return nil
}

func limit(body []byte) string {
	n := len(body)
	if n >= 100 {
		return string(body[n-100 : n])
	}
	return string(body)
}

func (u *Upstream) DelValidate() error {
	if u.Name == "" {
		return fmt.Errorf("wk_name not provided")
	}
	if u.Namespace == "" {
		return fmt.Errorf("namespace not provided")
	}
	if u.IP == "" {
		return fmt.Errorf("ip not provided")
	}
	if u.Port == "" {
		return fmt.Errorf("port not provided")
	}
	if u.Env == "" {
		return fmt.Errorf("env not provided")
	}
	if u.NginxGrp == "" {
		return fmt.Errorf("nginxgrp not provided")
	}
	return nil
}

func (u *Upstream) Del() error {
	if err := u.DelValidate(); err != nil {
		return fmt.Errorf("parameter for del validate err: %v", err)
	}
	return retry("del", u.del, 3)
}

func (u *Upstream) del() (err error) {
	resp, err := resty.SetRetryCount(3).
		//SetDebug(true).
		R().SetFormData(map[string]string{
		"wk_name":   u.Name,
		"namespace": u.Namespace,
		"ip_list":   u.IP,
		"port":      u.Port,
		"env":       u.Env,
		"nginx":     u.NginxGrp,
	}).
		Post(UpstreamDelAPI)

	if err != nil {
		return
	}
	log.Println("del resp: ", strings.Replace(limit(resp.Body()), "\n", "", -1))

	state, err := parseState(resp.Body())
	if err != nil {
		err = fmt.Errorf("upstream del err resp: %v", limit(resp.Body()))
		return
	}
	if state != true {
		err = fmt.Errorf("upstream del result failed resp: %v", limit(resp.Body()))
		return
	}
	return
}

// ChangeState change project specific ip state, remove item from nginx
// The logic may need to distinguish VM and docker
// We currently check based on docker first.
//
// Upstream will make it disabled ( need rethink?)
// Service recovery need human manual operation.
func ChangeState(endpoint, appname, namespace, state string) (ok bool, err error) {
	ip, port := endpoint2ip(endpoint)
	if ip == "" || port == "" {
		err = fmt.Errorf("ip: %v or port: %v not provided", ip, port)
		return
	}
	if appname == "" || namespace == "" || state == "" {
		err = fmt.Errorf("appname: %v or namespace: %v or state: %v not provided", appname, namespace, state)
		return
	}
	resp, err := resty.SetRetryCount(3).
		//SetDebug(true).
		R().SetFormData(map[string]string{
		"appname":   appname,
		"namespace": namespace,
		"ip":        ip,
		"port":      port,
		"state":     state, // int 1:up or 0:down
	}).
		Post(UpstreamnChangeAPI)

	if err != nil {
		return
	}

	log.Println("resp: ", limit(resp.Body()))

	result, err := parseState(resp.Body())
	if result != true {
		log.Println("ChangeState resp: ", limit(resp.Body()))
	}
	return result, err
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
