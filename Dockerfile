# Use an official Golang runtime as a parent image
FROM golang:alpine as build
WORKDIR /cli
COPY . /cli
RUN go get
RUN go build -o gitlab-runner-hadr

FROM alpine:3.19.0 as runtime
WORKDIR /cli
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=build cli/gitlab-runner-hadr /cli/gitlab-runner-hadr

# Define the command to run your application
ENTRYPOINT ["./gitlab-runner-hadr"]
