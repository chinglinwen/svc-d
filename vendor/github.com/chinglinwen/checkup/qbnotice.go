// qianbao need notice
package checkup

import (
	"log"
	"wen/svc-d/notice"
)

var (
	AlertAll        bool
	DefaultReceiver string
)

// it's depend on svc-d for flag setting
func Send(content, status, expire string) {
	reply, err := notice.Send(DefaultReceiver, content, status, expire)
	if err != nil {
		log.Printf("send alertall err: %v, resp: %v\n", err, reply)
	} else {
		log.Printf("send alertall ok, resp: %v\n", reply)
	}

}
