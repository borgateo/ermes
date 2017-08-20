FROM golang:1.8-alpine

RUN mkdir -p /go/src/github.com/borteo/ermes
ADD . /go/src/github.com/borteo/ermes/
WORKDIR /go/src/github.com/borteo/ermes

RUN apk add --no-cache git make

RUN go get ./...

RUN go build -ldflags "-X main.version=$(git rev-parse --short HEAD)" -o main .

CMD [ "/go/src/github.com/borteo/ermes/main" ]
