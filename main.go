package main

import (
	"flag"
	"os"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	"github.com/y3sh/go143/http"
	"github.com/y3sh/go143/instagram"
	"github.com/y3sh/go143/nytimes"
	"github.com/y3sh/go143/projects"
	"github.com/y3sh/go143/repository"
	"github.com/y3sh/go143/twitter"
)

var (
	serverPort  = flag.String("port", getEnv("PORT", "3000"), "Rest API Port, e.g. 3000")
	logLevelStr = flag.String("logLevel", getEnv("LOG_LEVEL", "debug"),
		"Log level: trace, debug, info, warn, error, fatal, panic")
)

func main() {
	flag.Parse()

	setupLogger(logLevelStr)

	redisPassword := os.Getenv("REDIS_PASSWORD")
	s3AccessKey := os.Getenv("S3_ACCESS_KEY")
	s3SecretKey := os.Getenv("S3_SECRET_KEY")
	nyTimesAPIKey := os.Getenv("NY_TIMES_API_KEY")

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
	nyTimesClient := nytimes.NewRestClient(nyTimesAPIKey)
	projectService := projects.NewProjectStoreService(redisRepository)

	chiRouter := chi.NewRouter()

	http.NewAPIRouter(chiRouter, tweetService, instagramUserService,
		nyTimesClient, projectService, s3Repository)

	restAPIServer, err := http.NewServer(http.Port(*serverPort))
	if err != nil {
		log.Fatalf("Failed to create api server. \n%+v\n", err)
	}

	log.Info("143 API starting . . .")

	err = restAPIServer.Start()
	if err != nil {
		log.Infof("Server shutdown.  \n%+v\n", err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func setupLogger(logLevelStr *string) {
	logLevel, err := log.ParseLevel(*logLevelStr)
	if err != nil {
		logLevel = log.TraceLevel
		log.Warn("Log level invalid or not provided, using trace level.")
	}

	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
		ForceColors:   true,
	})

	// Log filename and line number
	log.SetReportCaller(true)

	log.SetOutput(os.Stdout)
	log.SetLevel(logLevel)
}
