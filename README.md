# qframe-filter-id
qframe filter, which just passes input to the data or back channel


## main.go

When configured like this...

```go
cfgMap := map[string]string{
    "filter.test.send-back": "docker-event",
    "filter.test.send-data": "docker-event",
    "filter.test.inputs": "docker-events",
}
```
... an incoming `collector-docker-events` will be pushed to both channels.


```bash
$ go run main.go
2017/04/22 09:03:16 [II] Dispatch broadcast for Back, Data and Tick
2017/04/22 09:03:16 [II] Start id filter '%s' test
2017/04/22 09:03:18 Send message
#### Received message on Data-channel: docker-event
#### Received message on Back-channel: docker-event
```

## Development

```bash
$ docker run -ti --name qframe-collector-docker-events --rm -e SKIP_ENTRYPOINTS=1 \
           -v ${GOPATH}/src/github.com/qnib/qframe-filter-id:/usr/local/src/github.com/qnib/qframe-filter-id \
           -v ${GOPATH}/src/github.com/qnib/qframe-types:/usr/local/src/github.com/qnib/qframe-types \
           -v ${GOPATH}/src/github.com/qnib/qframe-utils:/usr/local/src/github.com/qnib/qframe-utils \
           -w /usr/local/src/github.com/qnib/qframe-filter-id \
            qnib/uplain-golang bash
root@835291194a41:/usr/local/src/github.com/qnib/qframe-filter-id# govendor update github.com/qnib/qframe-types github.com/qnib/qframe-utils
root@835291194a41:/usr/local/src/github.com/qnib/qframe-filter-id# govendor fetch +m
```
