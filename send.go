package main

import (
	"fmt"

	"github.com/sourcegraph/checkup"
)

type Wechat struct{}

func (s Wechat) Notify(results []checkup.Result) error {
	for _, result := range results {
		//if !result.Healthy {
		s.Send(result)
		//}
	}
	return nil
}

func (s Wechat) Send(result checkup.Result) error {
	fmt.Println("sended to wen", result)
	return nil
}
