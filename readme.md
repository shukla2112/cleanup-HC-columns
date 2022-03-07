```shell
ulimit -n 10240
curl https://api.honeycomb.io/1/columns/vector_multitenant_prod_filtered \
    -X GET \
    -H "X-Honeycomb-Team: <HC-token>" | jq . > hc-columns-202203.json
cat hc-columns-202203.json | grep -C 1 "detail.request.req_transform_ctx\|message.query.variables.where" | grep "id" | grep -v "message" | grep -v "hidden" | grep -v "detail.request.req_transform_ctx" | awk '{print $2}' | cut -d'"' -f2 > columnIds.txt
GODEBUG=netdns=go CGO_ENABLED=0 go run main.go
```
