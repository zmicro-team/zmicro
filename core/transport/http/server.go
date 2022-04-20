package http

import "net"

type Server struct {
}

func (s *Server) Start(l net.Listener) error {
	return nil
}

func (s *Server) Stop() error {
	return nil
}