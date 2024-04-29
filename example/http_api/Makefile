LDFLAGS := -ldflags "-s -w"
IMAGE := woshiaotian/httpapi:latest

.PHONY: build
build:
	env go build ${LDFLAGS} -o httpapi ./

.PHONY: image-multiple
image-multiple:
	docker buildx build -f ./Dockerfile \
		--platform linux/amd64,linux/arm64 --push -t ${IMAGE} .
