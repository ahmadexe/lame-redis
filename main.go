package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/ahmadexe/lame-redis/client"
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

func (s *Server) handleRawMessage(msg string) {
	cmd, err := parseCommand(msg)
	if err != nil {
		slog.Error("failed to parse command", "err", err)
		return
	}

	switch cmd.Name() {
	case CMD_SET:
		fmt.Printf("Set called, key: %s, value: %s\n", cmd.(SetCommand).Key, cmd.(SetCommand).Value)
	}

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

	msgChan := make(chan string)

	peer := NewPeer(conn, msgChan)

	s.addPeerChan <- peer

	slog.Info("new peer connected", "addr", conn.RemoteAddr())

	go s.handleReads(msgChan)

	if err := peer.readLoop(); err != nil {
		slog.Error("error while reading from peer", "remoteAddr", conn.RemoteAddr(), "err", err)
	}
}

func (s *Server) handleReads(msgChan chan string) {
	for msg := range msgChan {
		s.handleRawMessage(msg)
	}

	slog.Info("peer disconnected")
}

func main() {
	go func() {
		server := NewServer(Config{})
		err := server.Listen()
		if err != nil {
			slog.Error("error while starting server", "err", err)
		}
	}()
	time.Sleep(time.Second)
	client := client.NewClient("localhost:9090")
	if err := client.Set(context.Background(), "foo", "bar"); err != nil {
		slog.Error("failed to set", "err", err)
		panic("fail")
	}

	select {}
}
