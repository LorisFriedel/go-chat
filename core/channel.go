package core

import (
	"fmt"
	"net"
	"sync"

	"time"

	log "github.com/sirupsen/logrus"
)

type IChannel interface {
	Open() error
	Close() error
	Addr() net.Addr
}

type Channel struct {
	open     bool
	wg       sync.WaitGroup
	id       Identity
	registry IRegistry
	address  string
	port     int
	password string
	listener net.Listener
	msg      chan Message
	timeout  time.Duration
}

func NewChannel(address string, port int, password string, timeout int) *Channel {
	return &Channel{
		open:     false,
		id:       *NewIdentity(fmt.Sprintf("%s:%d", address, port)),
		registry: NewRegistry(),
		address:  address,
		port:     port,
		password: password,
		msg:      make(chan Message),
		timeout:  time.Duration(timeout) * time.Second,
	}
}

// Open make the channel listen for connection and handling received message
func (c *Channel) Open() error {
	err := c.listen()
	if err != nil {
		log.Errorln("Channel.Open: can't open channel connection")
		return err
	}

	log.Infof("Channel opened on %s\n", c.listener.Addr().String())
	return nil
}

func (c *Channel) listen() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", c.address, c.port))

	if err != nil {
		log.Errorf("Client.listen: can't listen on (%s)\n", fmt.Sprintf("%s:%d", c.address, c.port))
		return err
	}

	c.listener = listener
	c.open = true

	go c.handleJoin()
	go c.handleMsg()

	return nil
}

func (c *Channel) handleJoin() {
	c.wg.Add(1)
	for c.open {
		conn, err := c.listener.Accept()
		if err != nil {
			log.Errorf("Channel.handleJoin: connection error: %v", err)
			continue
		}
		go c.Join(conn)
		log.Infof("Channel.handleJoin: %v connected to channel", conn.LocalAddr().String())
	}
	c.wg.Done()
	log.Infoln("Channel.handleJoin: join handling is now inactive")
}

func (c *Channel) handleMsg() {
	for c.open {
		c.handle(<-c.msg)
	}
}

func (c *Channel) Close() (err error) {
	c.broadcastText("Channel closed by host.")

	// End infinite loop
	c.open = false

	// Close tcp connection
	err = c.listener.Close()

	// Wait for all loop to properly end
	c.wg.Wait()

	// Close client pipes
	c.registry.Foreach(func(id Identity, p *Pipe) {
		if p.IsOpen() {
			p.Close()
		}
	})

	close(c.msg)

	return
}

// Handle is used to handle message received from connected client
func (c *Channel) handle(msg Message) {
	// TODO Handle regarding message type ?
	c.broadcast(msg)
}

// Broadcast send the given message to every client connected to the channel
func (c *Channel) broadcast(msg Message) {
	log.Infof("Channel.broadcast: broadcasting message from: %s (%v)", msg.Sender, msg.Timestamp)

	c.registry.Foreach(func(id Identity, p *Pipe) {
		if p.IsOpen() {
			p.Write(msg) // TODO not ignore error ? I mean, who cares ?
		}
	})
}

// Join create a pipe between the channel and the given connection
// and start listening for client message, after executing the authentication procedure
// (the client must know the protocol)
func (c *Channel) Join(conn net.Conn) {
	// Open pipe to communicate with the client
	p := NewPipe(conn)

	// Listen for HELLO message, with client Identity
	msgHello, err := p.Read()
	if err != nil || msgHello.Type != HELLO {
		return
	}

	id := msgHello.Sender

	// If user is already in our registry, he is authenticated no need for password, send welcome back message
	if c.registry.Exists(id) { // TODO useless because every disconnected client is removed from the registry
		p.Write(*NewMsg(c.id, WELCOME_BACK)) // TODO handle write error
	} else if c.password != "" { // If not, ask for password if there is one
		p.Write(*NewMsg(c.id, PASSWORD_PLEASE)) // TODO handle write error
		msgPassword, err := p.Read()
		if err != nil || msgPassword.Type != PASSWORD {
			p.Write(*NewMsg(c.id, ERROR)) // TODO handle write error
			return
		}

		// If password match, OK
		if msgPassword.Text == c.password {
			p.Write(*NewMsg(c.id, WELCOME)) // TODO handle write error
		} else { // If not, close connection.
			p.Write(*NewMsg(c.id, WRONG_PASSWORD)) // TODO handle write error
			return
		}
	} else {
		p.Write(*NewMsg(c.id, WELCOME)) // TODO handle write error
	}

	c.registry.Push(id, p)
	c.broadcastText(fmt.Sprintf("%s joined the channel.", id.Name))
	log.Infof("Channel.Join: %s joined the channel (%s)", id.Name, id.Hash)

	// While client is connected
	for msg, err := p.Read(); p.IsOpen(); msg, err = p.Read() {
		if c.timeout > 0 && err.(net.Error).Timeout() {
			if p.IsOpen() {
				p.Write(*NewMsgSysChannel(c.id, fmt.Sprintf("..You sleepin', me kickin'")))
			}
			p.Close()
			c.broadcastText(fmt.Sprintf("%s has been inactive for %v and earned a nice and smooth KICK.", id.Name, c.timeout))
		}

		if err != nil {
			log.Errorf("Channel.Join: reading error while receiving client message: %v", err)
			continue
		}
		log.Infof("Channel.Join: received message from: %s (%v)", msg.Sender, msg.Timestamp)

		// Broadcast message
		c.msg <- msg

		if c.timeout > 0 {
			p.conn.SetReadDeadline(time.Now().Add(c.timeout * time.Second))
		}
	}

	if !c.open {
		return
	}

	// Here client is disconnected, pipe with him is closed
	c.registry.Pop(id)
	c.broadcastText(fmt.Sprintf("%s leaved the channel.", id.Name))
	log.Infof("Channel.Join: %s leaved the channel (%s)", id.Name, id.Hash)
}

func (c *Channel) broadcastText(text string) {
	c.msg <- *NewMsgSysChannel(c.id, text)
}

// Addr return the ip address of the channel
func (c *Channel) Addr() net.Addr {
	return c.listener.Addr()
}
