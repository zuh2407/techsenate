package network

import (
    "fmt"
    "net"
)

type Peer struct {
    ID   string
    Addr string
}

func NewPeer(id, addr string) *Peer {
    return &Peer{ID: id, Addr: addr}
}

func (p *Peer) SendMessage(msg string) error {
    conn, err := net.Dial("tcp", p.Addr)
    if err != nil {
        return err
    }
    defer conn.Close()

    fmt.Fprintf(conn, msg)
    return nil
}