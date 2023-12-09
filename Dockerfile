# Use an official Golang runtime as a parent image
FROM golang:alpine as build
WORKDIR /cli
COPY . /cli
RUN go get
RUN go build -o gitlab-runner-hadr

FROM alpine:3.19.0 as runtime
WORKDIR /cli
COPY --from=build cli/gitlab-runner-hadr /cli/gitlab-runner-hadr

# Define the command to run your application
ENTRYPOINT ["./gitlab-runner-hadr"]
