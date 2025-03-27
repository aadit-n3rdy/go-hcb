package main

import (
	"flag"
	"fmt"
	"net/rpc"
	//"github.com/aadit-n3rdy/hotCrossBuns/ds"
)

type Client struct {
	client *rpc.Client
}

// start the client and connect to the rpc server
func (c *Client) Start(port string) {
	client, err := rpc.Dial("tcp", "127.0.0.1:"+port)
	if err != nil {
		panic("Error connecting to replica")
	}

	c.client = client
}

// function to send command to replica
func (c *Client) SendCmd(cmd *string) {
	reply := false
	err := c.client.Call("Replica.ExecuteCommand", cmd, &reply)
	if err != nil {
		fmt.Println("Error sending command to replica: ", err)
	}
}

// function to commit changes
func (c *Client) Commit() {
	reply := false
	err := c.client.Call("Replica.Commit", &reply, &reply)
	if err != nil {
		fmt.Println("Error commiting to replica: ", err)
	}
}

// function to pull changes
func (c *Client) Pull() {
	reply := false
	err := c.client.Call("Replica.Pull", &reply, &reply)
	if err != nil {
		fmt.Println("Error pulling changes: ", err)
	}
}

func main() {

	// command line flags
	strPtr := flag.String("port", "8070", "Port number of replica")
	flag.Parse()

	c := new(Client)
	c.Start(*strPtr)

	inp := new(string)

	// read input from stdin and execute a client command based on it
	for {
		fmt.Scanf("%s", inp)

		switch *inp {
		case "stop":
			break
		case "commit":
			c.Commit()

		case "pull":
			c.Pull()
		default:
			c.SendCmd(inp)
		}
	}
}
