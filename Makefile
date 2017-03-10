
.PHONY: image tag build

GIT_DESCRIBE = $(shell git describe --tags)

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-X main.GitDescribe=$(GIT_DESCRIBE)'

image: build
	docker build -t zarnovican/go-helloworld .

tag:
	docker tag zarnovican/go-helloworld zarnovican/go-helloworld:$(GIT_DESCRIBE)
