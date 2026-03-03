package main

import (
	"log/slog"
	"net"
)

const defaultAddr = ":9090"

type Config struct {
	ListenAddr string
}

type Server struct {
	Config
	peers          map[*Peer]bool
	ln             net.Listener
	addPeerChan    chan *Peer
	removePeerChan chan *Peer
}

func NewServer(config Config) *Server {
	if config.ListenAddr == "" {
		config.ListenAddr = defaultAddr
	}

	return &Server{
		Config:         config,
		peers:          make(map[*Peer]bool),
		addPeerChan:    make(chan *Peer),
		removePeerChan: make(chan *Peer),
	}
}

func (s *Server) Listen() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	s.ln = ln

	go s.run()

	slog.Info("server started", "addr", s.ListenAddr)

	return s.acceptLoop()
}

func (s *Server) run() {
	for {
		select {
		case peer := <-s.addPeerChan:
			s.peers[peer] = true
		case peer := <-s.removePeerChan:
			delete(s.peers, peer)
		}
	}
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("error while accepting connection", "err", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	peer := NewPeer(conn)
	s.addPeerChan <- peer

	slog.Info("new peer connected", "addr", conn.RemoteAddr())

	peer.readLoop()
}

func main() {
	server := NewServer(Config{})
	err := server.Listen()
	if err != nil {
		slog.Error("error while starting server", "err", err)
	}
}
