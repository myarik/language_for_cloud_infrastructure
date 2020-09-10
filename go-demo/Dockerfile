FROM golang:alpine as build-backend

LABEL maintainer="Yaroslav Muravskyi <mail@myarik.com>"

# Set the Current Working Directory inside the container
WORKDIR /usr/src/app
# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN /usr/local/go/bin/go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY firstresponse/main.go .
# Build the Go app
RUN \
GOOS=linux GOARCH=amd64 /usr/local/go/bin/go build -o go-demo -ldflags "-X main.revision=0.0.1 -s -w" main.go

# (stage 2) package the binary into an Alpine container
FROM alpine:3.7 as app

COPY --from=build-backend /usr/src/app/go-demo /usr/bin/go-demo

CMD ["/usr/bin/go-demo"]