# build
FROM golang:1.21 as builder
ADD . /tmp/
WORKDIR /tmp

RUN go build -o /bin/httpapi -ldflags "-s -w" ./main.go


FROM woshiaotian/simple-base-image:v0.1.6

WORKDIR /data
COPY --from=builder /bin/httpapi /data/httpapi

CMD ["/data/httpapi"]