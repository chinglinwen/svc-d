package checkup

import (
	"encoding/json"
	"fmt"

	"gopkg.in/resty.v1"
)

type Qianbao struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Channel  string `json:"channel"`
	Webhook  string `json:"webhook"`
}

// Notify implements notifier interface
func (q Qianbao) Notify(results []Result) error {
	for _, result := range results {
		if !result.Healthy || result.Degraded {
			q.Send(result)
		}
	}
	return nil
}

// even error happen, continue next sends
func (q Qianbao) Send(r Result) error {
	/*
		var msg string
		if r.Message != "" {
			msg += ", msg: " + r.Message
		}
		if str := fmt.Sprintf("%v", r.Times); str != "" {
			msg += ", times: " + str
		}
		fmt.Printf("notify: send %v, %v, status: %v%v\n", r.Title, r.Endpoint, r.Status(), msg)
		//spew.Dump("result", r)
	*/
	b, err := json.Marshal(r)
	if err != nil {
		fmt.Println("notify marshal err:", err)
		return nil
	}
	_, err = resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(b).
		Post(q.Channel)
	if err != nil {
		fmt.Println("notify send err:", err)
		return nil
	}
	return nil
}
