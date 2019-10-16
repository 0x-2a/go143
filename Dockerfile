FROM golang:1.13.1-alpine3.10@sha256:2293e952c79b8b3a987e1e09d48b6aa403d703cef9a8fa316d30ba2918d37367 as builder

RUN apk update && apk upgrade && apk add --no-cache \
 ca-certificates \
 && rm -rf /var/cache/apk/*

# Create unprivileged appuser
RUN adduser -D -g '' appuser

WORKDIR /go143-build

# Copy source from machine dir to workdir
COPY . .

ENV GO111MODULE="on"

# ldflags -w -s omits the symbol table, debug information and the DWARF table
RUN go get -u ./...
RUN GOOS=linux GOOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" \
 	-o bin/go143 ./main.go

# Use a small image for runtime
FROM scratch

WORKDIR /bin

# Copy appuser creds
COPY --from=builder /etc/passwd /etc/passwd

# Copy certs for ssl requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy built binary
COPY --from=builder /go143-build/bin/go143 .

# Use an unprivileged user.
USER appuser

ENTRYPOINT ["/bin/go143"]
