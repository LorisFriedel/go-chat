package core

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"sync"

	log "github.com/sirupsen/logrus"
)

type IPipe interface {
	Read() (Message, error)
	Write(Message) error
	Close() error
}

type Pipe struct {
	conn    net.Conn
	open    bool
	wg      sync.WaitGroup
	decoder *json.Decoder
	encoder *json.Encoder
}

var ErrPipeClosed error = errors.New("pipe is closed")

func NewPipe(conn net.Conn) *Pipe {
	pipe := &Pipe{
		conn:    conn,
		open:    true,
		decoder: json.NewDecoder(conn),
		encoder: json.NewEncoder(conn),
	}

	return pipe
}

func (p *Pipe) Read() (msg Message, err error) {
	err = p.decoder.Decode(&msg)
	if err == io.EOF {
		p.Close()
		err = ErrPipeClosed
	}
	return
}

func (p *Pipe) Write(msg Message) (err error) {
	err = p.encoder.Encode(msg)
	if err != nil {
		log.Errorf("Pipe.write: message encoding error: %v\n", err)
	}
	return
}

func (p *Pipe) Close() error {
	if !p.open {
		return ErrPipeClosed
	}
	p.open = false

	log.Infoln("Pipe.Close: closing pipe connection")
	err := p.conn.Close() // break the infinite loop
	if err != nil {
		log.Errorln("Pipe.Close: net.Conn not properly closed")
		return err
	}

	log.Infoln("Pipe.Close: waiting for pipe to properly close")
	p.wg.Wait()

	return nil
}

func (p *Pipe) IsOpen() bool {
	return p.open
}
