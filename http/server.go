package http

import (
	"fmt"
	"net/http"

	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
)

type Port string

type Config struct {
	Port string
}

type Server struct {
	Port Port
}

func NewServer(port Port) (*Server, error) {
	server := &Server{
		Port: port,
	}

	return server, nil
}

func (s *Server) Start() error {
	err := s.listen()
	if err != nil {
		return errors.Trace(err)
	}

	return nil
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
