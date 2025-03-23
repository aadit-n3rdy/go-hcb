package types

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type Config struct {
	NodeAddrs      []string // ip:port of all nodes
	SelfIndex      int      // index of this node in NodeAddrs
	FaultTolerance int      // value of "f" (maximum number of Byzantine nodes)
}

func (config *Config) Read() {
	configFilePtr := flag.String("config", "", "The path to the JSON config file")

	flag.Parse()

	f, err := os.Open(*configFilePtr)
	if err != nil {
		fmt.Println("Error opening config file:", err)
		return
	}
	defer f.Close()
	configDec := json.NewDecoder(f)

	var configMap map[string]any
	err = configDec.Decode(&configMap)
	if err != nil {
		fmt.Println("Invalid config file:", err)
		return
	}

	selfIndexStr, ok := os.LookupEnv("HCB_SELF")
	if !ok {
		fmt.Println("Missing HCB_SELF env var")
		return
	}
	var selfIndex int
	fmt.Sscan(selfIndexStr, &selfIndex)

	config.SelfIndex = selfIndex
	v, ok := configMap["nodes"]
	if !ok {
		fmt.Println("Config is missing list \"nodes\"")
		return
	}
	nodeAddrs, ok := v.([]string)
	if !ok {
		fmt.Println("List of nodes was not list of string")
		return
	}
	config.NodeAddrs = nodeAddrs
	config.FaultTolerance = (len(nodeAddrs) - 1) / 3
}
