LDFLAGS = -ldflags "-s -w"

.PHONY: build
build:
	env CGO_ENABLED=0 go build ${LDFLAGS} -o autotest ./
xgo:
	xgo --targets=linux/amd64,darwin/amd64,darwin/arm64 -out=./autotest  .
