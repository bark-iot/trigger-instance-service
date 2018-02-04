FROM golang:1.9.3-alpine

RUN apk update && apk add curl git

# Install dep
RUN curl -L https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 > /bin/dep &&\
    chmod +x /bin/dep

WORKDIR /go/src/app

COPY . .

# Install dependencies
RUN dep ensure