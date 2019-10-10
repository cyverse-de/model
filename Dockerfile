FROM golang:1.9-alpine

RUN apk update && apk add git

RUN go get github.com/jstemmer/go-junit-report
RUN go get github.com/golang/dep/cmd/dep

COPY . /go/src/github.com/cyverse-de/model

RUN cd /go/src/github.com/cyverse-de/model && dep ensure

CMD go test -v github.com/cyverse-de/model | tee /dev/stderr | go-junit-report
