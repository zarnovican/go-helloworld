
.PHONY: image tag

image:
	docker build -t zarnovican/go-helloworld .

tag:
	docker tag zarnovican/go-helloworld zarnovican/go-helloworld:$(shell git describe --tags)
