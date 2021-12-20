package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	go143http "github.com/y3sh/go143/http"
	"github.com/y3sh/go143/instagram"
	"github.com/y3sh/go143/nytimes"
	"github.com/y3sh/go143/polygon"
	"github.com/y3sh/go143/projects"
	"github.com/y3sh/go143/proxyURL"
	"github.com/y3sh/go143/repository"
	"github.com/y3sh/go143/twitter"
)

const (
	dialTimeout       = 5 * time.Second
	handshakeTimeout  = 5 * time.Second
	responseTimeout   = time.Second * 10
	DebugTSFormat     = "2006-01-02 03:04:05PM MST"
	longestFileLength = 28
)

func main() {
	flag.Parse()

	logLevel := os.Getenv("LOG_LEVEL")
	SetupLogger(logLevel)

	redisPassword := os.Getenv("REDIS_PASSWORD")
	s3AccessKey := os.Getenv("S3_ACCESS_KEY")
	s3SecretKey := os.Getenv("S3_SECRET_KEY")
	nyTimesAPIKey := os.Getenv("NY_TIMES_API_KEY")
	polygonAPIKey := os.Getenv("POLYGON_API_KEY")
	googleBooksAPIKey := os.Getenv("GOOGLE_BOOKS_API_KEY")

	serverHost := os.Getenv("HOST")
	if serverHost == "" {
		serverHost = "localhost"
	}

	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = "8080"
	}
	hostAddress := fmt.Sprintf("%s:%s", serverHost, serverPort)

	redisRepository := repository.NewRedisRepository()
	err := redisRepository.Connect(redisPassword)
	if err != nil {
		log.Fatalf("Failed to connect to redis. \n%+v\n", err)
	}

	s3Repository, err := repository.NewS3Repository(s3AccessKey, s3SecretKey)
	if err != nil {
		log.Fatalf("Failed to connect to s3. \n%+v\n", err)
	}

	err = redisRepository.Connect(redisPassword)
	if err != nil {
		log.Fatalf("Failed to connect to redis. \n%+v\n", err)
	}

	tweetService := twitter.NewTweetService()
	instagramUserService := instagram.NewUserService()
	nyTimesClient := nytimes.NewRestClient(nyTimesAPIKey, googleBooksAPIKey, GetHTTPClient())
	polygonClient := polygon.NewRestClient(polygonAPIKey, GetHTTPClient())
	proxyClient := proxyURL.NewProxyClient(GetHTTPClient())
	projectService := projects.NewProjectStoreService(redisRepository)

	chiRouter := chi.NewRouter()

	go143http.NewAPIRouter(chiRouter, tweetService, instagramUserService,
		nyTimesClient, polygonClient, proxyClient, projectService, s3Repository)

	log.Infof("REST API starting on %s . . .", hostAddress)
	err = http.ListenAndServe(hostAddress, chiRouter)
	if err != nil {
		log.Fatalf("HTTP Server Exited:  \n%s", errors.ErrorStack(err))
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func SetupLogger(logLevelStr string) {
	if logLevelStr == "" {
		logLevelStr = "trace"
	}

	logLevel, err := log.ParseLevel(logLevelStr)
	if err != nil {
		logLevel = log.TraceLevel
	}

	log.SetFormatter(&log.TextFormatter{
		DisableColors:          false,
		FullTimestamp:          true,
		ForceColors:            true,
		TimestampFormat:        DebugTSFormat,
		PadLevelText:           true,
		DisableLevelTruncation: true,
		DisableSorting:         true,
		CallerPrettyfier: func(frame *runtime.Frame) (string, string) {
			fileStr := frame.File

			idx := strings.LastIndex(fileStr, "/")
			if idx > -1 {
				fileStr = fileStr[idx+1:]
			}

			fileLine := fmt.Sprintf(" %s:%d", fileStr, frame.Line)

			for len(fileLine) < longestFileLength {
				fileLine += " "
			}

			return "", fileLine
		},
	})

	log.SetReportCaller(true)
	log.SetOutput(os.Stdout)
	log.SetLevel(logLevel)

	log.Infof("Logger started with %s level.", log.GetLevel())
}

func GetHTTPClient() *http.Client {
	httpTransport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: dialTimeout,
		}).Dial,
		TLSHandshakeTimeout: handshakeTimeout,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{
		Timeout:   responseTimeout,
		Transport: httpTransport,
	}

	return httpClient
}
