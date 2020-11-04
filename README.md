# go143

A fast API written in go for COS 143 Web Dev

## Dev Setup

* [Go](https://golang.org/dl/) version should be at 1.13.1+ to support go modules
* [docker](https://docs.docker.com/docker-for-mac/) should be at least 18.0X

```sh
# No need to be inside go path, modules are enabled
git clone git@github.com:y3sh/go143.git
cd go143

# Pull down the libs
go get -u ./... 
```

## Running via source

Compile and run:

```sh
go build -o ./bin/go143 .
./bin/go143
```

Flags include `port` and `logLevel` (panic, fatal, error, warn, info, debug, trace)
```sh
./bind/go143 --port=8080 --logLevel=trace
```

Test the API at [http://localhost:3000/](http://localhost:3000/)

## Start a Redis Instance

docker run \
-p 6379:6379 \
-v redisData:/data \
--name redis \
--restart on-failure \
-d redis:6.0.9-alpine redis-server --appendonly yes  --requirepass "REDIS_PASSWORD_HERE"


## Running via Docker

```sh
docker build -f Dockerfile -t go143:1.0.0 .

docker run --rm -ti -p 8080:8080 -e REDIS_PASSWORD="REDIS_PASSWORD_HERE" go143:1.0.0 --port=8080 --logLevel=trace

docker stop go143:1.0.0
```

## Running PROD via Docker

```sh
sudo docker run -d --restart on-failure -p 3000:8080 go143:1.0.0 --port=8080 --logLevel=info
```

## Linting

1. Install https://github.com/golangci/golangci-lint
2. Add your $GOPATH/bin/golangci-lint to your path
3. Run `./goLint.sh`

