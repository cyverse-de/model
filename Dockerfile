FROM golang:1.6-alpine

RUN apk update && apk add git

RUN go get github.com/spf13/viper
RUN go get github.com/cyverse-de/configurate

COPY . /go/src/github.com/cyverse-de/model

CMD ["go", "test", "github.com/cyverse-de/model"]
