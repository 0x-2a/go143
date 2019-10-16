package http

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/juju/errors"

	log "github.com/sirupsen/logrus"
)

const (
	SiteRoot  = "/"
	APIRoot   = "/api"
	TweetsURI = "/api/v1/tweets"
)

type Port string

type Config struct {
	Port string
}

type Server struct {
	APIRouter *APIRouter
	Port      Port
}

func NewServer(router *APIRouter, port Port) (*Server, error) {
	server := &Server{
		APIRouter: router,
		Port:      port,
	}

	server.bindRoutes()

	return server, nil
}

func (s *Server) Start() error {
	err := s.listen()
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (s *Server) bindRoutes() {
	chiRouter := chi.NewRouter()

	corsConfig := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           100, // Maximum value not ignored by any of major browsers
	})

	chiRouter.Use(corsConfig.Handler)

	chiRouter.Route(SiteRoot, s.APIRouter.HandleRoot)
	chiRouter.Route(APIRoot, s.APIRouter.HandleRoot)
	chiRouter.Route(TweetsURI, s.APIRouter.HandleTweets)

	http.Handle(SiteRoot, chiRouter)
}

func (s *Server) listen() error {
	portAddr := fmt.Sprintf(":%s", s.Port)
	log.Infof("Server listening on %s", portAddr)

	err := http.ListenAndServe(portAddr, http.DefaultServeMux)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}
