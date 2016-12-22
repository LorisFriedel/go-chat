package core

import "fmt"

// TODO interface ?
type Client struct {
	name          string
	channel       *Channel
	knownChannels []*Channel
	ownChannels   []*Channel

	// Each channel assign a unique ID for each client
}

func NewClient(name string) *Client {
	return &Client{name, nil, []*Channel{}, []*Channel{}}
}

func (c *Client) Name() string {
	return c.name
}

func (c *Client) CurrentChannel() *Channel {
	return c.channel
}

func (c *Client) SendMessage(msg string) error {
	fmt.Println(c.channel, " : ", msg)
	// TODO channel.sendMessage
	return nil
}

// TODO send message [to a channel]

// TODO disconnect [from a channel]
// TODO change current channel, get current satus
// TODO get server status (if owner/or not, info is different)
// TODO change server password (if owner)
// TODO

func (c *Client) String() string {
	return c.name
}
