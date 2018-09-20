// notice send wechat msg to the people.
package notice

import (
	"flag"
	"strings"

	"gopkg.in/resty.v1"
)

var (
	noticeBase = flag.String("noticeurl", "http://noti.wk.qianbao-inc.com", "notice api base url")
)

// if no need cache, leave expire empty, then will send every time.
func Send(receiver, message, status, expire string) (reply string, err error) {
	r := strings.NewReplacer("\"", " ", "{", "", "}", "")
	message = r.Replace(message)

	resp, e := resty.SetRetryCount(3).R().
		SetQueryParams(map[string]string{
			"user":    receiver,
			"content": message,
			"status":  status,
			"expire":  expire,
		}).
		Get(*noticeBase)

	if e != nil {
		err = e
		return
	}
	reply = string(resp.Body())
	return
}
