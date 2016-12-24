package core

import (
	"fmt"
	"github.com/golang/glog"
	"log"
	"net"
)

type IChannel interface {
	Open() error
	Close() error
	Broadcast(Message)
	Join(net.Conn)
	Addr() net.Addr
}

type Channel struct {
	clients  []*Pipe
	joins    chan net.Conn
	incoming chan Message
	outgoing chan Message
	address  string
	port     int
	passwd   string
	listener net.Listener
}

// TODO accept only net.Addr ????!!!!!
func newChannel(address string, port int, passwd string) *Channel {
	return &Channel{
		clients:  make([]*Pipe, 0),
		joins:    make(chan net.Conn),
		incoming: make(chan Message),
		outgoing: make(chan Message),
		address:  address,
		port:     port,
		passwd:   passwd,
	}
}

func (c *Channel) Open() error {
	err := c.listen()
	if err != nil {
		glog.Errorln("Channel.Open: can't open channel connection")
		return err
	}

	glog.Infof("Channel oppened on %s\n", c.listener.Addr().String())
	return nil
}

func (c *Channel) listen() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", c.address, c.port))

	if err != nil {
		glog.Errorf("Client.listen: can't listen on (%s)\n", fmt.Sprintf("%s:%d", c.address, c.port))
		return err
	}

	c.listener = listener

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Connection error: %v", err)
				continue
			}
			// TODO check closed channel ?
			c.joins <- conn
			// TODO password authentication if not known..
			// TODO !!! i.e. accept connection but start auth procedure, if not successful close connection.
		}
	}()

	go func() {
		for {
			// TODO check closed channel ?
			select {
			case msg := <-c.incoming:
				glog.Infof("received message from: %s (%v)", msg.Sender, msg.Timestamp)
				c.Broadcast(msg)
			case conn := <-c.joins:
				glog.Infof("%v joined the channel", conn.LocalAddr().String())
				c.Join(conn)
			}
		}
	}()

	return nil
}

func (c *Channel) Close() error {
	// TODO Send bye message to all listener
	close(c.joins)

	for _, client := range c.clients {
		client.Close()
	}

	close(c.incoming)
	close(c.outgoing)

	return c.listener.Close() // TODO call first instead of last ?
}

func (c *Channel) Broadcast(msg Message) {
	glog.Infof("broadcasting message from: %s (%v)", msg.Sender, msg.Timestamp)

	for _, client := range c.clients {
		// TODO check closed channel ? no in the other function
		client.outgoing <- msg
	}
}

func (c *Channel) Join(conn net.Conn) {
	client := newPipe(conn)
	client.Open()
	c.clients = append(c.clients, client)

	go func() {
		//defer client.Close()

		// TODO check closed channel ? OK ! ?
		for msg, more := <-client.incoming; more; msg, more = <-client.incoming {
			c.incoming <- msg
		}
	}()
}

// TODO function remove/delete/bye

func (c *Channel) Addr() net.Addr {
	return c.listener.Addr()
}

/*******************/

type KnownChan struct {
	name    string
	address string
	port    int
	passwd  string // TODO ?
}

// TODO le password, une fois qu'il est rentré dans l'app, est hashé menu et on l'envoie comme ça
func newKnownChan(name, address string, port int, passwd string) *KnownChan {
	return &KnownChan{
		name:    name,
		address: address,
		port:    port,
		passwd:  passwd,
	}
}

// TODO fonction de getConnection() pour pouvoir s'y reconnecter (avec authentiication à chaque fois + stockage de MDP ou bien c'est
// TODO le serveur qui se rappelle de nous ?????????)

// c'est une pipeline a double sens
/*
Il faut qu'il soit sur écoute permannente et puisse envoyé aussi
Le client a son channel courant enregistré (observer) en tant que recepteur

// Il faut la procédure de connexion/authentifaction bien défini et a part

On aura deux routine, une qui envoie au channel courant dès qu'on met dans le pipeline
et une qui recoi dès que quelqu'un nous envoi un message

// TODO check les channels qui ont pas déjà la même addres / port

// TODO max client !!!
*/