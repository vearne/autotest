LDFLAGS = -ldflags "-s -w"

.PHONY: build
build:
	env go build ${LDFLAGS} -o autotest ./
xgo:
	xgo -out=./autotest  .
