FROM scratch

ADD go-helloworld /

ENTRYPOINT ["/go-helloworld"]
