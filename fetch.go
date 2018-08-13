// periodically fetch checks info
package main

// call upstream api, transformation to  checkers
// how to know if it's http or tcp?

// does it support setting for every check? info come from platform?

// sync all or, sync one?

// check info constantly change? checker need update?
// we usually check the old info, if things changed, resync?

// 5 minutes interval

import (
	"encoding/json"
	"time"

	"github.com/chinglinwen/checkup"
	"github.com/chinglinwen/log"
	"github.com/go-resty/resty"
)

func start() {
	log.Println("started fetch in the background")

	var (
		ticker = time.NewTicker(10 * time.Second)
		//stop   *time.Ticker
		now  time.Time
		prev time.Time
	)

	for ; true; now = <-ticker.C {
		log.Printf("starting fetch at %v, last time elapsed %v\n", now.Format(time.UnixDate), now.Sub(prev))
		// empty it for new fetch
		conf.Checkers = []checkup.Checker{}

		projects, err := fetch(*upstreamAPI)
		if err != nil {
			log.Println("fetch error", err)
			continue
		}

		convertToCheck(projects)
		log.Println("fetch ok")

		err = conf.CheckAndStore()
		if err != nil {
			log.Println("check error", err)
			continue
		}
		prev = now

		if *checkonetime {
			ticker.Stop()
			goto EXIT
		}
	}
EXIT:
}

// later may adopt setting for projects
func convertToCheck(projects Projects) {
	var envcnt, enabledcnt, dockercnt, checkcnt int
	for _, p := range projects {
		if p.ENV != *env {
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
		if *testproject != "" && p.Project != *testproject {
			continue
		}
		checkcnt++

		check := checkup.HTTPChecker{
			Name:     p.Project,
			URL:      "http://" + p.IP + ":" + p.Port,
			Attempts: 5,
			UpStatus: 452, //above 500 consider error
		}
		//log.Printf("project %v added\n", p.Project)
		conf.Checkers = append(conf.Checkers, check)
	}
	log.Printf("got %v projects, %v enabled, %v is docker, %v to be check\n",
		envcnt, enabledcnt, dockercnt, checkcnt)
	if *testproject != "" {
		log.Println("will only run for testing project: ", *testproject)
	}
	conf.Save()
}

type Projects []Prjoect

type Prjoect struct {
	Project  string `json:"project",omitempty"`
	Enabled  string `json:"enabled",omitempty"`
	ENV      string `json:"env",omitempty"`
	IDC      string `json:"idc",omitempty"`
	IP       string `json:"ip",omitempty"`
	Port     string `json:"port",omitempty"`
	IsDocker string `json:"is_docker",omitempty"`
}

func fetch(url string) (p Projects, err error) {
	if url == "" {
		err = json.Unmarshal([]byte(testitems), &p)
		return
	}
	resp, e := resty.R().Get(url)
	if e != nil {
		err = e
		return
	}
	err = json.Unmarshal(resp.Body(), &p)
	return
}

var testitems = `
[
  {
    "enabled": "1", 
    "env": "test", 
    "idc": "BJ-SH", 
    "ip": "fs.qianbao-inc.com", 
    "is_docker": "1", 
    "port": "80", 
    "project": "ops_test"
  }, 
  {
    "enabled": "1", 
    "env": "test", 
    "idc": "BJ-SH", 
    "ip": "104.16.25.88", 
    "is_docker": "1", 
    "port": "80", 
    "project": "ip.cn"
  }, 
  {
    "enabled": "1", 
    "env": "test", 
    "idc": "BJ-SH", 
    "ip": "example.com", 
    "is_docker": "1", 
    "port": "80", 
    "project": "example"
  }
]
`
