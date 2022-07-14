package ctx

import (
	"net"
	"x-game/message"
)

type Socket interface {
	Read()
	Write()
}

type socketModel struct {
	conn  net.Conn
	close chan struct{}
	Msg   chan message.Message
}

func newSocketModel(c net.Conn) *socketModel {
	m := &socketModel{conn: c, close: make(chan struct{})}
	return m
}

func (m *socketModel) run() {

}
