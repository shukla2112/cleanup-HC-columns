```shell
ulimit -n 10240
cat ~/work/Hasura/hc-columns-202203.json | grep -C 1 "detail.request.req_transform_ctx\|message.query.variables.where" | grep "id" | grep -v "message" | grep -v "hidden" | grep -v "detail.request.req_transform_ctx" | awk '{print $2}' | cut -d'"' -f2 > columnIds.txt
GODEBUG=netdns=go CGO_ENABLED=0 go run main.go
```
