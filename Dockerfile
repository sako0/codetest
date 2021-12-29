FROM golang:1.16

ENV DB_HOST "172.18.0.1"

COPY app /go/app
COPY go.mod /go
COPY go.sum /go

WORKDIR /go

ENV GOPATH=/go/app

CMD go run ./app/main.go