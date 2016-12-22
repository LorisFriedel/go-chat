package core

import "fmt"

// TODO interface ?
type Client struct {
	Name string
	Channel *Channel // pointer ?
	// Known channel
	// Name
	// Own channel
	// Current channel
	// Each channel assign a unique ID for each client
}

func NewClient(name string) *Client {
	return &Client{name, nil} /// TODO channel etc..
}

func (c *Client) SendMessage(msg string, channel *Channel) error {
	fmt.Println(channel, " : ", msg)
	// TODO
	return nil
}

// TODO send message [to a channel], disconnect [from a channel], change current channel, get current satus
// TODO get server status (if owner/or not, info is different), change server password (if owner)
// TODO

func (c *Client) String() string {
	return c.Name
}