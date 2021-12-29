FROM golang:1.16

COPY app /go/app
COPY go.mod /go
COPY go.sum /go

WORKDIR /go

CMD go run ./app/main.go