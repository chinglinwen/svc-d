package check

import (
	"encoding/json"
	"flag"
	"fmt"

	"gopkg.in/resty.v1"
)

var (
	ProjectCenterBase   = flag.String("center-api", "http://ops.qianbao-inc.com/api", "project center all config api base url")
	ProjectCenterAllAPI = *ProjectCenterBase + "/get_detect_projects"
	ProjectCenterAPI    = *ProjectCenterBase + "/get_detect_project"
	projectDockerInfo   = *ProjectCenterBase + "/project_docker_info"
)

type ProjectCheck struct {
	ProjectName    string `json:"project_name"`
	Enabled        string `json:"enabled"`
	AttemptSpacing int    `json:"attempt_spacing"`
	Attempts       int    `json:"attempts"`
	MustContain    string `json:"must_contain"`
	StatusCode     int    `json:"status_code"`
	Timeout        int    `json:"timeout"`
	Type           string `json:"type"` // http or tcp
	URI            string `json:"uri"`

	AutoDisable   bool   `json:"auto_disable"` // auto disable or not
	AlertReceiver string `json:"alert_receiver"`
}

type ConfigBody struct {
	Code int          `json:"code,omitempty"`
	Msg  string       `json:"message,omitempty"`
	Data ProjectCheck `json:"data,omitempty"`
}

func FetchConfigByK8sName(name string) (config ProjectCheck, err error) {
	appname, err := GetProjectName(name)
	if err != nil {
		return
	}
	return FetchConfig(appname)
}

// the name need to be underscore one, say ops_test
func FetchConfig(name string) (config ProjectCheck, err error) {
	var p ConfigBody
	resp, e := resty.SetRetryCount(3).
		//SetDebug(true).
		R().Get(ProjectCenterAPI + "/" + name)
	if e != nil {
		err = e
		return
	}
	p, err = decodeConfig(resp.Body())
	if err != nil {
		return
	}

	if p.Code != 0 {
		err = fmt.Errorf("%v", p.Msg)
		return
	}
	config = p.Data
	return
}

type ConfigsBody struct {
	Code int            `json:"code,omitempty"`
	Msg  string         `json:"message,omitempty"`
	Data []ProjectCheck `json:"data,omitempty"`
}

type ProjectChecks map[string]ProjectCheck

func FetchConfigs() (configs ProjectChecks, err error) {
	var p ConfigsBody
	resp, e := resty.SetRetryCount(3).
		//SetDebug(true).
		R().Get(ProjectCenterAllAPI)
	if e != nil {
		err = e
		return
	}

	p, err = decodeConfigs(resp.Body())
	if err != nil {
		return nil, err
	}

	if p.Code != 0 {
		return nil, fmt.Errorf("%v", p.Msg)
	}

	configs = make(ProjectChecks, len(p.Data))
	for _, v := range p.Data {
		configs[v.ProjectName] = v
	}
	return
}

func decodeConfig(body []byte) (p ConfigBody, err error) {
	err = json.Unmarshal(body, &p)
	return
}

func decodeConfigs(body []byte) (p ConfigsBody, err error) {
	err = json.Unmarshal(body, &p)
	return
}

func GetProjectName(wkname string) (name string, err error) {
	resp, err := resty.SetRetryCount(3).
		//SetDebug(true).
		R().
		SetFormData(map[string]string{
			"project_cluster_name_d": wkname,
		}).
		Post(projectDockerInfo)

	if err != nil {
		return
	}
	type body struct {
		Projectname string `json:"project_name"`
	}
	var b body
	err = json.Unmarshal(resp.Body(), &b)
	if err != nil {
		return
	}
	return b.Projectname, nil
}
