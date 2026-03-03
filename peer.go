package main

import (
	"net"
)

type Peer struct {
	conn net.Conn
	msgChan chan string
}

func NewPeer(conn net.Conn, msgChan chan string) *Peer {
	return &Peer{conn: conn, msgChan: msgChan}
}

func (p *Peer) readLoop() error {
	buf := make([]byte, 1024)
	defer close(p.msgChan)

	for {
		n, err := p.conn.Read(buf)
		if err != nil {
			return err
		}

		msg := string(buf[:n])
		p.msgChan <- msg
	}
}
