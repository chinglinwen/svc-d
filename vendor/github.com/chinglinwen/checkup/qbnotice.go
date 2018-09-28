// qianbao need notice
package checkup

import (
	"fmt"
	"wen/svc-d/notice"

	"github.com/chinglinwen/checkup/cache"
	"github.com/chinglinwen/log"
)

var (
	AlertAll        bool
	DefaultReceiver string
)

// it's depend on svc-d for flag setting
func Send(content, status, expire string) {
	// avoid repeat sent
	user := DefaultReceiver
	if d, found := cache.Get(user, content, status); found {
		e := fmt.Sprintf("user: %v, content: %v not expired in %v, skip send\n",
			user, content, d.Format("15:04:05"))
		log.Debug.Printf(e)
		return
	}

	// set cache
	log.Printf("user %v, %v, status: %v, expire set to %v\n", user, content, status, expire)
	cache.Set(user, content, status, expire)

	reply, err := notice.Send(user, content, status, expire)
	if err != nil {
		log.Printf("send alertall err: %v, resp: %v\n", err, reply)
	} else {
		log.Printf("send alertall ok, resp: %v\n", reply)
	}

}
