FROM golang:1.12.5-alpine3.9 as build-stage

ENV GO111MODULE off
RUN apk --no-cache add git && go get github.com/54mch4n/kamachat
ENV GO111MODULE on

WORKDIR /go/src/github.com/54mch4n/kamachat
RUN GOOS=linux GOARCH=amd64 go build -o ./bin/ojichat

FROM alpine:latest as exec-stage
COPY --from=build-stage /go/src/github.com/54mch4n/kamachat/bin/kamachat /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/kamachat"]
