package ctx

import (
	"common_module/g"
	"common_module/message"
	"github.com/golang/protobuf/proto"
	"io"
	"log"
	"net"
	"reflect"
)

type Conner interface {
	Start()
	Stop()
	ReadMsg()
	WriteMsg()
}

type socketConn struct {
	conn  net.Conn
	In    chan *message.Message
	Out   chan []byte
	close chan struct{}
}

func NewConner(c net.Conn) *socketConn {
	m := &socketConn{
		conn:  c,
		In:    make(chan *message.Message, g.MsgLength),
		Out:   make(chan []byte, g.MsgLength),
		close: make(chan struct{}),
	}
	return m
}

func (m *socketConn) Start() {
	go m.msgHandle()
	go m.WriteMsg()
	m.ReadMsg()
}

func (m *socketConn) Stop() {
	m.waitMsgHandle()
	close(m.close)
	m.conn.Close()
}

func (m *socketConn) ReadMsg() {
	for {
		buf := make([]byte, g.ReadWriteMaxLength)
		n, err := m.conn.Read(buf)
		if err != nil || err != io.EOF {
			log.Printf("read msg err,error is [%v]", err.Error())
			break
		}
		if n < g.ReadWriteMinLength {
			log.Printf("msg not match condition!")
			continue
		}
		msg := message.BytesToMsg(buf)
		m.In <- msg
	}
}

func (m *socketConn) msgHandle() {
	for {
		select {
		case msg := <-m.In:
			m.Handle(msg)
		case <-m.close:
			break
		}
	}
}

func (m *socketConn) WriteMsg() {
	for {
		select {
		case bytes := <-m.Out:
			_, err := m.conn.Write(bytes)
			if err != nil {
				log.Printf("write msg error[%v]", err.Error())
			}
		case <-m.close:
			break
		}

	}
}

func (m *socketConn) Handle(msg *message.Message) {
	hm, ok := message.GetHandlerModel(msg.MsgId)
	if !ok {
		return
	}
	data := reflect.New(hm.T.Elem()).Interface()
	if err := proto.Unmarshal(msg.Data, data); err != nil {
		log.Printf("Handle err,error is [%v]", err.Error())
		return
	}
	hm.H(data, nil)
}

func (m *socketConn) waitMsgHandle() {
	close(m.In)
	for {
		msg, ok := <-m.In
		if !ok {
			break
		}
		m.Handle(msg)
	}
}
