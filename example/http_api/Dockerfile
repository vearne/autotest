# build
FROM golang:1.21 as builder
ADD ../../ $GOPATH/src/github.com/vearne/autotest
WORKDIR $GOPATH/src/github.com/vearne/autotest/example/http_api

RUN go build -o /bin/httpapi -ldflags "-s -w"


FROM woshiaotian/simple-base-image:v0.1.6

WORKDIR /data
COPY --from=builder /bin/httpapi /data/httpapi

CMD ["/data/httpapi"]