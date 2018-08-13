#!/bin/sh
GOOS=linux go build
curl -s fs.qianbao-inc.com/t/uploadapi -F file=@svc-d -F truncate=yes
