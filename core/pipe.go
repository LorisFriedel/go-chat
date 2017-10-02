package core

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"sync"

	"bufio"

	"fmt"

	log "github.com/sirupsen/logrus"
)

type Interceptor func(p *Pipe, msg *Message)

type IPipe interface {
	Read() (Message, error)
	Write(Message) error
	Close() error
}

type Pipe struct {
	conn     net.Conn
	open     bool
	wg       sync.WaitGroup
	decoder  *json.Decoder
	encoder  *json.Encoder
	muR      sync.Mutex
	muW      sync.Mutex
	raw      bool
	wIntrcpt Interceptor
}

var ErrPipeClosed error = errors.New("pipe is closed")

func NewPipe(conn net.Conn) *Pipe {
	pipe := &Pipe{
		conn:    conn,
		open:    true,
		decoder: json.NewDecoder(conn),
		encoder: json.NewEncoder(conn),
		muR:     sync.Mutex{},
		muW:     sync.Mutex{},
		raw:     false,
	}

	return pipe
}

func (p *Pipe) Read() (msg Message, err error) {
	p.muR.Lock()
	defer p.muR.Unlock()

	if p.raw {
		var msgStr string
		msgStr, err = bufio.NewReader(p.conn).ReadString('\n')
		msg = *NewMsgText(Identity{}, msgStr)
	} else {
		err = p.decoder.Decode(&msg)
		if err == io.EOF {
			p.Close()
			err = ErrPipeClosed
		}
	}

	if err != nil {
		log.Errorf("Pipe.Read: %v <-- %v: message decoding error: %v", p.conn.LocalAddr(), p.conn.RemoteAddr(), err)
	}

	return
}

func (p *Pipe) Write(msg Message) (err error) {
	p.muW.Lock()
	defer p.muW.Unlock()

	if p.wIntrcpt != nil {
		p.wIntrcpt(p, &msg)
	}

	if p.raw {
		fmt.Fprintln(p.conn, msg.Text)
	} else {
		err = p.encoder.Encode(msg)
	}

	if err != nil {
		log.Errorf("Pipe.Write: %v --> %v: message encoding error: %v", p.conn.LocalAddr(), p.conn.RemoteAddr(), err)
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
