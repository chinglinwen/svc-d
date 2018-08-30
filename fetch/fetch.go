// periodically fetch checks info
package fetch

// call upstream api, transformation to  checkers
// how to know if it's http or tcp?

// does it support setting for every check? info come from platform?

// sync all or, sync one?

// check info constantly change? checker need update?
// we usually check the old info, if things changed, resync?

// 5 minutes interval

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"wen/hook-api/upstream"

	"github.com/chinglinwen/checkup"
	"github.com/chinglinwen/log"
	"github.com/go-resty/resty"
)

var (
	UpstreamAllAPI    = upstream.UpstreamAllAPI
	UpstreamSingleAPI = upstream.UpstreamSingleAPI
)

var (
	// limit env for fetched projects
	Env = "qa" // qa,pre,pro
)

// for get all projects struct
type Projects []Project

type Project struct {
	Name      string `json:"project,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Enabled   string `json:"enabled,omitempty"`
	Env       string `json:"env,omitempty"`
	IDC       string `json:"idc,omitempty"`
	IP        string `json:"ip,omitempty"`
	Port      string `json:"port,omitempty"`
	IsDocker  string `json:"is_docker,omitempty"`
}

// for single project
type AProject struct {
	Name      string `json:"upstream_name,omitempty"`
	Namespace string `json:"project,omitempty"`
	IP        string `json:"upstream_ip,omitempty"`
	Port      string `json:"upstream_port,omitempty"`
	State     int    `json:"upstream_state,omitempty"`
	NginxGrp  string `json:"upstream_group,omitempty"`
	OpsCode   int    `json:"ops_code,omitempty"`
}

// all ips for a project
type AProjectIps struct {
	Env  string
	Name string
	IPs  []AProject
}

// many aproject mapp one project, so we pass config instead
func (p *AProject) GetCheck(config ProjectCheck) (check checkup.Checker) {
	if config.Type == "http" {
		c := checkup.HTTPChecker{
			Name: p.Name + "#" + p.Namespace + "#" + "http", // notify api will use it
			URL:  "http://" + p.IP + ":" + p.Port,
		}
		if config.StatusCode == 0 {
			c.UpStatus = 452 //above 500 consider error
		} else {
			c.UpStatus = config.StatusCode
		}
		if config.Attempts != 0 {
			c.Attempts = config.Attempts
		}
		if config.MustContain != "" {
			c.MustContain = config.MustContain
		}
		check = c
	} else {
		c := checkup.TCPChecker{
			Name: p.Name + "#" + p.Namespace + "#" + "tcp", // notify api will use it
			URL:  p.IP + ":" + p.Port,
		}
		if config.Timeout != 0 {
			c.Timeout = time.Duration(config.Timeout) * time.Second
		}
		check = c
	}
	return
}

func (a AProjectIps) Check() error {
	if len(a.IPs) == 0 {
		return fmt.Errorf("no ip found for project: %v", a.Name)
	}
	checks := make([]checkup.Checker, 0)
	config, err := FetchConfig(a.Name)
	if err != nil {
		// just log
		log.Printf("fetch config for %v, err: %v\n", a.Name, err)
	}

	for _, v := range a.IPs {
		check := v.GetCheck(config)
		checks = append(checks, check)
	}

	c := checkup.Checkup{
		Checkers: checks,
	}
	r, err := c.Check()
	if err != nil {
		return err
	}
	for _, v := range r {
		if !v.Healthy {
			return fmt.Errorf("%v, healthy: %v", v.Endpoint, v.Healthy)
		}
	}
	return nil
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

/*
func FetchAppName(env, name string) (p AProjectIps, err error) {
	return fetch(env, "appname", name)
}

func Fetch(env, name string) (p AProjectIps, err error) {
	return fetch(env, "wk_name", name)
}
*/
func Fetch(env, name string) (p AProjectIps, err error) {
	p.Env = env
	p.Name = name
	resp, e := resty.SetRetryCount(3).R().
		SetQueryParam("env", env).
		SetQueryParam("appname", name).
		Get(UpstreamSingleAPI)
	if e != nil {
		err = e
		return
	}
	//fmt.Println(string(resp.Body()))

	// unmarshal error, probably error result body
	err = json.Unmarshal(resp.Body(), &p.IPs)
	return
}

func Fetchs() (p Projects, err error) {
	resp, e := resty.SetRetryCount(3).R().Get(UpstreamAllAPI)
	if e != nil {
		err = e
		return
	}
	err = json.Unmarshal(resp.Body(), &p)
	return
}

func ConvertToCheck(projects Projects, configs ProjectChecks) []checkup.Checker {
	return ConvertToCheckWithTest(projects, configs, "")
}

// later may adopt setting for projects
func ConvertToCheckWithTest(projects Projects, configs ProjectChecks, testproject string) []checkup.Checker {
	var envcnt, enabledcnt, dockercnt, checkcnt int

	checks := make([]checkup.Checker, 0)

	for _, p := range projects {
		if p.Env != Env {
			continue
		}
		envcnt++
		if p.Enabled == "0" {
			continue
		}
		enabledcnt++

		if p.IsDocker != "1" {
			continue
		}
		dockercnt++

		// test first
		if testproject != "" && p.Name != testproject {
			continue
		}
		checkcnt++

		if p.Name == "" || p.Namespace == "" {
			log.Println("got invalid project", p)
			continue
		}

		var check checkup.Checker
		if configs[p.Name].Type == "http" {
			c := checkup.HTTPChecker{
				Name: p.Name + "#" + p.Namespace + "#" + "http", // notify api will use it
				URL:  "http://" + p.IP + ":" + p.Port,
			}

			if configs[p.Name].StatusCode == 0 {
				c.UpStatus = 452 //452, //above 500 consider error
			} else {
				c.UpStatus = configs[p.Name].StatusCode
			}
			if configs[p.Name].Attempts != 0 {
				c.Attempts = configs[p.Name].Attempts
			}
			if configs[p.Name].MustContain != "" {
				c.MustContain = configs[p.Name].MustContain
			}
			check = c
		} else {
			c := checkup.TCPChecker{
				Name: p.Name + "#" + p.Namespace + "#" + "tcp", // notify api will use it
				URL:  p.IP + ":" + p.Port,
			}
			if configs[p.Name].Timeout != 0 {
				c.Timeout = time.Duration(configs[p.Name].Timeout) * time.Second
			}
			check = c
		}

		//spew.Dump("check", checks)
		//log.Printf("project %v added\n", p.Project)
		checks = append(checks, check)
	}
	log.Printf("got %v projects, %v enabled, %v is docker, %v to be check\n",
		envcnt, enabledcnt, dockercnt, checkcnt)

	return checks

}
