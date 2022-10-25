FROM golang:1.19.2-alpine3.16
WORKDIR /build

# Fetch dependencies
COPY go.mod go.sum ./
RUN go mod download
COPY main.go main.go  
RUN go build main.go
cmd ["./main"]