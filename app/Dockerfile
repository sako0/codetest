FROM golang:1.16

COPY . /go/app

WORKDIR /go

RUN mv /go/app/go.mod /go
RUN mv /go/app/go.sum /go

ENV GOPATH=/go/app

CMD go run ./app/main.go