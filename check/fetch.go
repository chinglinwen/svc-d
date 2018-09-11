// periodically fetch checks info
package check

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

	"github.com/chinglinwen/checkup"
	"github.com/chinglinwen/log"
	"github.com/go-resty/resty"
)

var (
	// limit env for fetched projects
	Env        = "qa" // qa,pre,pro
	DockerOnly bool

	UpstreamBase      string
	UpstreamAllAPI    string
	UpstreamSingleAPI string
)

func Init(UpstreamBase string) {
	UpstreamAllAPI = UpstreamBase + "/get_upstream_all_instance/"
	UpstreamSingleAPI = UpstreamBase + "/get_nginx_all/"
}

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

	var name string
	if config.Type == "http" {
		name = p.Name + "#" + p.Namespace + "#" + "http"
	} else {
		name = p.Name + "#" + p.Namespace + "#" + "tcp"
	}
	check = GetCheckWithConfig(name, p.IP, p.Port, config)

	return
}

func (a AProjectIps) Check() (r []checkup.Result, err error) {
	if len(a.IPs) == 0 {
		err = fmt.Errorf("no ip found for project: %v", a.Name)
		return
	}
	checks := make([]checkup.Checker, 0)
	config, err := FetchConfig(a.Name)
	if err != nil {
		// just log, later return err to the platform?
		err = fmt.Errorf("fetch config for %v, err: %v", a.Name, err)
		log.Println(err)
		//return
	}

	for _, v := range a.IPs {
		check := v.GetCheck(config)
		checks = append(checks, check)
	}

	c := checkup.Checkup{
		Checkers: checks,
	}
	log.Debug.Println("c definition", c)

	return c.Check()
}

func Endpoint2ip(e string) (ip, port string) {
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

// Merge upstream info and project center's check config.
func ConvertToCheck(projects Projects, configs ProjectChecks) []checkup.Checker {
	var envcnt, enabledcnt, dockercnt, checkcnt int

	checks := make([]checkup.Checker, 0)

	for _, p := range projects {
		if p.Env != Env {
			continue
		}
		envcnt++
		if p.Enabled != "1" {
			// upstream not enable
			log.Debug.Printf("%v disabled by upstream\n", p.Name)
			continue
		}
		if configs[p.Name].Enabled != "on" {
			// project center not enable
			log.Debug.Printf("%v disabled by platform\n", p.Name)
			continue
		}
		enabledcnt++

		if DockerOnly {
			if p.IsDocker != "1" {
				continue
			}
			dockercnt++
		}
		checkcnt++

		if p.Name == "" || p.Namespace == "" {
			log.Println("got invalid project", p)
			continue
		}

		var name string
		config := configs[p.Name]
		if config.Type == "http" {
			name = p.Name + "#" + p.Namespace + "#" + "http"
		} else {
			name = p.Name + "#" + p.Namespace + "#" + "tcp"
		}
		check := GetCheckWithConfig(name, p.IP, p.Port, config)

		//spew.Dump("check", checks)
		//log.Printf("project %v added\n", p.Project)
		checks = append(checks, check)
	}
	log.Printf("got %v projects, %v enabled, %v is docker, %v to be check\n",
		envcnt, enabledcnt, dockercnt, checkcnt)

	return checks

}
