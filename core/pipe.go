package core

import (
	"encoding/json"
	"github.com/golang/glog"
	"net"
	"sync"
	"io"
)

type Pipe struct {
	conn     net.Conn
	open     bool
	wg       sync.WaitGroup
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
	p.wg.Add(1)
	for p.open {
		var msg Message
		err := p.decoder.Decode(&msg)
		if err != nil {
			if err == io.EOF {
				p.Close()
				break
			} else {
				glog.Errorf("Pipe.read: message decoding error: %v", err)
				continue
			}
		}
		p.incoming <- msg
	}
	p.wg.Done()
	glog.Infof("Pipe.read: stop reading process for channel %v\n", p.conn.RemoteAddr())
}

func (p *Pipe) write() {
	for msg := range p.outgoing {
		glog.Infof("Pipe.write: sending message from %v to %v (%v)\n",
			p.conn.LocalAddr(), p.conn.RemoteAddr(), msg.Timestamp)
		err := p.encoder.Encode(msg)
		if err != nil {
			glog.Errorf("Pipe.write: message encoding error: %v\n", err)
		}
	}
	glog.Infof("Pipe.write: stop reading process for channel %v\n", p.conn.LocalAddr())
}

func (p *Pipe) Open() {
	p.open = true
	go p.read()
	go p.write()
}

func (p *Pipe) Close() (err error) {
	p.open = false

	glog.Infoln("Pipe.Close: closing outgoing channel")
	close(p.outgoing) // break the infinite loop

	glog.Infoln("Pipe.Close: closing pipe connection")
	err = p.conn.Close() // break the infinite loop
	if err != nil {
		glog.Errorln("Pipe.Close: net.Conn not properly closed")
	}

	glog.Infoln("Pipe.Close: waiting for pipe to properly close")
	p.wg.Wait()

	glog.Infoln("Pipe.Close: closing incoming channel")
	close(p.incoming)

	return
}

func (p *Pipe) IsOpen() bool {
	return p.open // TODO better check ?
}
