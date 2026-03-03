package main

import "net"

type Config struct {
	ListenAddr string

}

type Server struct {
	Config
	ln net.Listener
}

func NewServer(config Config) *Server {
	return &Server{
		Config: config,
	}
}

func main()  {
	
}