package fcgiserver

import (
	"fmt"
	"net"
	"net/http"
	"net/http/fcgi"

	"github.com/s2ks/fcgiserver/logger"
)

type Handler interface {
	Setup(string) error
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type FcgiServer struct {
	ServeMux *http.ServeMux

	Address  string
	Port     string
	Protocol string
}

func (s *FcgiServer) Register(path string, h Handler) {
	mux := s.ServeMux
	err := h.Setup(path)

	if err != nil {
		logger.Fatal(err)
	}

	mux.Handle(path, h)
	logger.Infof("Registered handler for %s", path)
}

func New(address string, port string, protocol string) (*FcgiServer, error) {
	var s *FcgiServer

	s = new(FcgiServer)

	s.ServeMux = http.NewServeMux()

	s.Address = address
	s.Port = port
	s.Protocol = protocol

	return s, nil
}

func (s *FcgiServer) Serve() error {
	listener, err := net.Listen(s.Protocol, fmt.Sprintf("%s:%s", s.Address, s.Port))

	defer listener.Close()

	if err != nil {
		return err
	}

	handler := s.ServeMux

	return fcgi.Serve(listener, handler)
}
