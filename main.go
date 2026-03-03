package main

import "net"

const defaultAddr = ":9090"

type Config struct {
	ListenAddr string

}

type Server struct {
	Config
	ln net.Listener
}

func NewServer(config Config) *Server {
	if config.ListenAddr == "" {
		config.ListenAddr = defaultAddr
	}

	return &Server{
		Config: config,
	}
}

func (s *Server) Listen() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	s.ln = ln
	return nil
}

func main()  {
	
}