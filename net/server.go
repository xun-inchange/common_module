package ctx

import (
	"log"
	"net"
	"sync"
)

type Server struct {
	Addr       string
	Listener   net.Listener
	SocketConn map[Conner]struct{}
}

func NewServer(addr string) *Server {
	m := &Server{}
	m.Addr = addr
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panicf("listen addr[%v] err,error is [%v]", addr, err.Error())
	}
	m.Listener = l
	m.SocketConn = make(map[Conner]struct{})
	return m
}

func (m *Server) Run() {
	for {
		conn, err := m.Listener.Accept()
		if err != nil {
			log.Printf("accept socket err,error is [%v]", err.Error())
			continue
		}
		conner := NewConner(conn)
		m.StoreSocketConn(conner)
		go conner.Start()
	}
}

func (m *Server) StoreSocketConn(c Conner) {
	m.SocketConn[c] = struct{}{}
}

func (m *Server) Close() {
	err := m.Listener.Close()
	if err != nil {

	}
	var wg sync.WaitGroup
	for c, _ := range m.SocketConn {
		wg.Add(1)
		go c.Stop()
		wg.Done()
	}
	wg.Wait()
}
