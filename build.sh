#!/bin/sh
GOOS=linux go build
tar --exclude=statuspage/js/config.js -czf svc-d.tar.gz svc-d statuspage
curl -s fs.qianbao-inc.com/k8s/soft/uploadapi -F file=@svc-d.tar.gz -F truncate=yes
cksum ./svc-d
