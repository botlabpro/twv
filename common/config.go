package common

import (
	"encoding/json"
	"fmt"
	"github.com/botlabpro/twv/node"
	"log"
	"os"
)

type Config struct {
	Nodes []*node.BaseNode `json:"nodes"`
}

func LoadConfig(file string) *Config {
	f, err := os.ReadFile(file)
	if err != nil {
		log.Println(err)
	}
	var data Config
	if err := json.Unmarshal(f, &data); err != nil {
		fmt.Println(err.Error())
		panic("no config file")
	}
	return &data
}
