LDFLAGS = -ldflags "-s -w"

.PHONY: build
build:
	env CGO_ENABLED=0 go build ${LDFLAGS} -o autotest ./
xgo:
	xgo -out=./autotest  .
