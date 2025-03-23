package main

import (
	"fmt"
	"math/rand"
	"net/rpc"

	"github.com/aadit-n3rdy/hotCrossBuns/types"
)

type Client struct {
	srvConns []*rpc.Client
	config   *Config
}

func (c *Client) Init(conf *Config) {
	c.config = conf
	c.srvConns = make([]*rpc.Client, len(conf.NodeAddrs))
	var err error
	for i, addr := range c.config.NodeAddrs {
		c.srvConns[i], err = rpc.Dial("tcp", addr)
		if err != nil {
			fmt.Println("Error connecting to ", addr, ":", err)
		}
	}
}

func (c *Client) SendMessage(topic string, msg string) {
	idxs := rand.Perm(len(c.srvConns))
	done := 0
	for _, i := range idxs {
		sent := false
		err := c.srvConns[i].Call("Server.Publish", types.NewMessage(topic, msg), &sent)
		if err != nil {
			fmt.Println("Error sending:", err)
		}
		if err == nil && sent {
			done += 1
		}
		if done > c.config.FaultTolerance {
			break
		}
	}
}

func main() {
	config := Config{}
	config.Read()
}
