FROM golang:1.13-alpine3.11 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.con/Xarakiri/L0

COPY go.mod go.sum ./
COPY vendor vendor
COPY util util
COPY cache cache
COPY event event
COPY db db
COPY schema schema
COPY query-service query-service
COPY pusher-service pusher-service

RUN GO111MODULE=on go install -mod vendor ./...

FROM alpine:3.11
WORKDIR /usr/bin
COPY --from=build /go/bin .
