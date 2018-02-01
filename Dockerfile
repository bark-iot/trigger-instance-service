FROM golang:1.9.2

# Install bee command line tool
RUN go get github.com/beego/bee

WORKDIR /go/src/app

COPY . .