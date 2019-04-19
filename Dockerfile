# Start from an image with the latest version of Go installed
FROM golang as builder
LABEL stage=intermediate

WORKDIR /workspace

# First we copy only the dependency list, and download those. (Improves build speed if code changes but not dependencies)
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the rest of the code and build the service
ADD . /workspace

RUN go test -v -race .

RUN GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -o /go/bin/service -ldflags "-X main.VERSION=$(git describe --tags --always) -X main.BUILD_TIME=`date -u +%Y%m%d.%H%M%S`"

# Now we use the output binary to make a lean image, where the ssh key and the source code are no longer present
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/service .
COPY --from=builder /workspace/config/config.template.json ./config/config.json

EXPOSE 8080
VOLUME /static
VOLUME /config
VOLUME /output

ENTRYPOINT ["/service"]