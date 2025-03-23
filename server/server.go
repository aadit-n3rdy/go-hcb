package main

import (
	"hotCrossBuns/server/types"
	"net/rpc"
)

type Server struct {
	index   int
	clients []*rpc.Client
	config  *types.Config
}

func main() {
	var config types.Config
	config.Read()
}
