## Build
FROM golang:1.20.10-alpine3.18 AS buildenv

RUN apk add --no-cache build-base

ADD go.mod go.sum /

RUN go mod download

WORKDIR /app

ADD . .

ENV GO111MODULE=on

ENV CGO_ENABLED=1

RUN  go build -o main cmd/main.go

## Deploy
FROM alpine

WORKDIR /

COPY --from=buildenv  /app/ /

EXPOSE 3000

CMD ["/main"]