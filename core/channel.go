package core

import (
	"fmt"
	"net"
	"sync"

	"time"

	"bufio"

	"strings"

	"encoding/json"

	"io"

	"github.com/LorisFriedel/go-chat/format"
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
	msg      chan *Message
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
		msg:      make(chan *Message),
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

	log.Infof("Channel opened on %s", c.listener.Addr().String())
	return nil
}

func (c *Channel) listen() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", c.address, c.port))

	if err != nil {
		log.Errorf("Client.listen: can't listen on (%s)", fmt.Sprintf("%s:%d", c.address, c.port))
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
func (c *Channel) handle(msg *Message) {
	// TODO Handle regarding message type ?
	c.broadcast(msg)
}

// Broadcast send the given message to every client connected to the channel
func (c *Channel) broadcast(msg *Message) {
	log.Infof("Channel.broadcast: broadcasting message from: %s (%v)", msg.Sender, msg.Timestamp)

	c.registry.Foreach(func(id Identity, p *Pipe) {
		if p.IsOpen() {
			if p.raw && id == msg.Sender {
				return
			}

			go p.Write(*msg) // TODO not ignore error ? I mean, who cares ?
		}
	})
}

// Join create a pipe between the channel and the given connection
// and start listening for client message, after executing the authentication procedure
// (the client must know the protocol)
func (c *Channel) Join(conn net.Conn) {
	// Open pipe to communicate with the client
	p := NewPipe(conn)

	fmt.Fprintln(p.conn, "Welcome stranger! What's your name?")
	msgHelloStr, err := bufio.NewReader(conn).ReadString('\n')

	var msgHello Message
	err = json.Unmarshal([]byte(msgHelloStr), &msgHello)

	// Listen for HELLO message, with client Identity
	if err != nil {
		msgHelloStr = strings.TrimSpace(msgHelloStr)
		if msgHelloStr == "" {
			log.Error("Channel.Join: not a JSON message and empty message")
			return
		}

		c.newRawClient(p, msgHelloStr)
		return
	} else if msgHello.Type != HELLO {
		log.Error("Channel.Join: wrong hello message type")
	}

	id := msgHello.Sender

	// If user is already in our registry, he is authenticated no need for password, send welcome back message
	if c.registry.Exists(id) { // TODO useless because every disconnected client is removed from the registry
		p.Write(*NewMsg(c.id, WELCOME_BACK)) // TODO handle write error
	} else if c.password != "" { // If not, ask for password if there is one
		p.Write(*NewMsg(c.id, PASSWORD_PLEASE)) // TODO handle write error
		msgPassword, err := p.Read()
		if err != nil || msgPassword.Type != PASSWORD {
			log.Errorf("Channel.Join: expected password message, got %v or error (%v)", msgPassword.Type, err)
			p.Write(*NewMsg(c.id, ERROR)) // TODO handle write error
			return
		}

		// If password match, OK
		if msgPassword.Text == c.password {
			p.Write(*NewMsg(c.id, WELCOME)) // TODO handle write error
		} else { // If not, close connection.
			p.Write(*NewMsg(c.id, WRONG_PASSWORD)) // TODO handle write error
			log.Errorf("Channel.Join: client %v entered a wrong password", c.id, err)
			return
		}
	} else {
		p.Write(*NewMsg(c.id, WELCOME)) // TODO handle write error
	}

	c.registry.Push(id, p)
	c.broadcastText(fmt.Sprintf("%s joined the channel.", id.Name))
	log.Infof("Channel.Join: %s joined the channel (%s)", id.Name, id.Hash)

	if c.timeout.Seconds() > 0 {
		p.conn.SetReadDeadline(time.Now().Add(c.timeout))
	}

	// While client is connected
	for msg, err := p.Read(); p.IsOpen(); msg, err = p.Read() {
		if err != nil {
			errTout, ok := err.(net.Error)
			if c.timeout.Seconds() > 0 && ok && errTout.Timeout() {
				if p.IsOpen() {
					p.Write(*NewMsgSysChannel(c.id, fmt.Sprintf("..You sleep, I kick!")))
				}
				p.Close()
				c.broadcastText(fmt.Sprintf("%s has been inactive for %v and earned a nice and smooth KICK.", id.Name, c.timeout))
			}

			log.Errorf("Channel.Join: reading error while receiving client message: %v", err)
			continue
		}
		log.Infof("Channel.Join: received message from: %s (%v)", msg.Sender, msg.Timestamp)

		// Broadcast message
		c.msg <- &msg

		if c.timeout > 0 {
			p.conn.SetReadDeadline(time.Now().Add(c.timeout))
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

func (c *Channel) newRawClient(p *Pipe, userName string) {
	log.Infof("Channel.newRawClient: starting new raw client for %v, %v", userName, p)

	if c.password != "" { // If not, ask for password if there is one
		fmt.Fprintln(p.conn, "Password?")

		msgPassword, err := bufio.NewReader(p.conn).ReadString('\n')
		msgPassword = strings.TrimSpace(msgPassword)

		if err != nil {
			log.Errorf("Channel.newRawClient: reading error while receiving password: %v", err)
			return
		}

		if msgPassword != c.password {
			fmt.Fprintln(p.conn, "Wrong password. Bye.")
			p.Close()
			return
		}
	}

	// create ID for this man and define a way to broadcast differently regarding user
	id := *NewIdentityFromConn(userName, p.conn)

	p.raw = true
	p.wIntrcpt = func(p *Pipe, msg *Message) {
		color := colorMap[msg.Type]
		if msg.Sender == id {
			color = format.LIGHT_GREEN
		}
		msg.Text = format.Msg(msg.Sender.Name, strings.TrimSpace(msg.Text), msg.Timestamp, color)
		// Induce a 'format' package dependency, but it is the way to allow a raw client to have formatted messages
	}

	c.registry.Push(id, p)
	p.Write(*NewMsgSysChannel(id, "Welcome on board "+userName+"!"))

	c.broadcastText(fmt.Sprintf("%s joined the channel.", id.Name))
	log.Infof("Channel.Join: %s joined the channel (%s)", id.Name, id.Hash)

	if c.timeout.Seconds() > 0 {
		p.conn.SetReadDeadline(time.Now().Add(c.timeout))
	}

	// While client is connected
	for msg, err := p.Read(); p.IsOpen(); msg, err = p.Read() {
		if err != nil {
			errTout, ok := err.(net.Error)
			if c.timeout.Seconds() > 0 && ok && errTout.Timeout() {
				if p.IsOpen() {
					p.Write(*NewMsgSysChannel(c.id, fmt.Sprintf("..You sleep, I kick!")))
				}
				p.Close()
				c.broadcastText(fmt.Sprintf("%s has been inactive for %v and earned a nice and smooth KICK.", id.Name, c.timeout))
			}

			if err == io.EOF {
				break
			}

			log.Errorf("Channel.Join: reading error while receiving client message: %v", err)
			continue
		}
		log.Infof("Channel.Join: received message from: %s (%v)", msg.Sender, msg.Timestamp)

		// Format message
		msg.Sender = id

		// Broadcast message
		c.msg <- &msg

		if c.timeout > 0 {
			p.conn.SetReadDeadline(time.Now().Add(c.timeout))
		}
	}

	// TODO timeout not working

	if !c.open {
		return
	}

	// Here client is disconnected, pipe with him is closed
	c.registry.Pop(id)
	c.broadcastText(fmt.Sprintf("%s leaved the channel.", id.Name))
	log.Infof("Channel.Join: %s leaved the channel (%s)", id.Name, id.Hash)
}

var colorMap map[TMsg]int = map[TMsg]int{
	TEXT:        format.LIGHT_BLUE,
	SYS_CHANNEL: format.LIGHT_RED,
	SYS_CLIENT:  format.RED,
}

func (c *Channel) broadcastText(text string) {
	c.msg <- NewMsgSysChannel(c.id, text)
}

// Addr return the ip address of the channel
func (c *Channel) Addr() net.Addr {
	return c.listener.Addr()
}

func (c *Channel) String() string {
	return fmt.Sprintf("name: %s, address: %v, password: %s, timeout: %v", c.id.Name, c.Addr(), c.password, c.timeout)
}
