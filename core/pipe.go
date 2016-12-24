package core

import (
	"encoding/json"
	"github.com/golang/glog"
	"log"
	"net"
)

type Pipe struct {
	conn     net.Conn
	open     bool
	incoming chan Message
	outgoing chan Message
	decoder  *json.Decoder
	encoder  *json.Encoder
}

func newPipe(conn net.Conn) *Pipe {
	pipe := &Pipe{
		conn:     conn,
		open:     false,
		incoming: make(chan Message),
		outgoing: make(chan Message),
		decoder:  json.NewDecoder(conn),
		encoder:  json.NewEncoder(conn),
	}

	// TODO maybe open here ?
	return pipe
}

func (p *Pipe) read() {
	for p.open {
		var msg Message
		err := p.decoder.Decode(&msg)
		if err != nil {
			log.Printf("message decoding error: %v", err)
			continue
		}
		// TODO check closed channel
		p.incoming <- msg
	}
	glog.Infoln("stop reading process for channel %v", p.conn.RemoteAddr())
}

func (p *Pipe) write() {
	// TODO check closed channel ?
	for msg := range p.outgoing {
		if !p.open {
			break
		}

		glog.Infof("Pipe.write: sending message from %v to %v (%v)\n",
			p.conn.LocalAddr(), p.conn.RemoteAddr(), msg.Timestamp)
		err := p.encoder.Encode(msg)
		if err != nil {
			glog.Errorf("message encoding error: %v\n", err)
		}
	}
	glog.Infoln("stop reading process for channel %v", p.conn.LocalAddr())
}

func (p *Pipe) Open() {
	p.open = true
	go p.read()
	go p.write()
}

func (p *Pipe) Close() error {
	// TODO proper close, safer close
	p.open = false
	close(p.outgoing)
	close(p.incoming)

	err := p.conn.Close()
	if err != nil {
		glog.Errorln("Pipe.Close: net.Conn not properly closed")
	}
	return err
}

func (p *Pipe) IsOpen() bool {
	return p.open // TODO proper check if pipe is really open ?
}
