# build
FROM golang:1.21 as builder

RUN mkdir -p /tmp/http_api
COPY go.mod /tmp/http_api/go.mod
COPY main.go /tmp/http_api/main.go

WORKDIR /tmp/http_api
RUN go build -o /bin/httpapi -ldflags "-s -w" /tmp/http_api


FROM woshiaotian/simple-base-image:v0.1.8

WORKDIR /data
COPY --from=builder /bin/httpapi /data/httpapi

CMD ["/data/httpapi"]