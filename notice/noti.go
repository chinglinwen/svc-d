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

func Send(receiver, message string) (reply string, err error) {
	r := strings.NewReplacer("\"", " ", "{", "", "}", "")
	message = r.Replace(message)

	resp, e := resty.SetRetryCount(3).R().
		SetQueryParams(map[string]string{
			"user":    receiver,
			"content": message,
		}).
		Get(*noticeBase)

	if e != nil {
		err = e
		return
	}
	reply = string(resp.Body())
	return
}
