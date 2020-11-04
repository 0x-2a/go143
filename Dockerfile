FROM golang:1.15.0-alpine3.12@sha256:73182a0a24a1534e31ad9cc9e3a4bb46bb030a883b26eda0a87060f679b83607 as builder

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
RUN GOOS=linux GOOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -mod vendor \
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
