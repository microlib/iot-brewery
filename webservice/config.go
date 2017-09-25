package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

/**
 * The json config will always be in the form of parent.key:value pair
 * the reasoning here is that it is easy to maintain and use
 * and also if required can be migtrated to a key value store such as redis
 *
 * Don't dig it - then feel welcome to change it to your hearts content - knock yourself out
 *
 **/

type Config struct {
	Level   string `json:"level"`
	Basedir string `json:"base_dir"`
	Port    string `json:"port"`
	Cache   string `json:"cache"`
	Servers []Server
}

type Server struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"pwd"`
}

// As the logger can only be configured after we read the config
// I make use of the stdout for error logging
func (cfg Config) Init(fileName string) Config {
	start := time.Now()
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("%s \x1b[1;31m[%s] \x1b[0m : %s \n", "ERROR", err)
		os.Exit(0)
	}

	err = json.Unmarshal(file, &cfg)
	if err != nil {
		fmt.Printf("%s \x1b[1;31m[%s] \x1b[0m : %s \n", "ERROR", err)
		os.Exit(0)
	}

	fmt.Printf("%s \x1b[1;34m [%s] \x1b[0m  : %s \n", start.Format("2006/01/02 03:04:05"), "INFO", "Config data read")
	return cfg
}

func (cfg Config) GetServer(name string) (Server, error) {

	var index int = -1
	for i, server := range cfg.Servers {
		if server.Name == name {
			index = i
			break
		}
	}

	if index == -1 {
		return Server{}, errors.New("Config could not find server by name")
	} else {
		return cfg.Servers[index], nil
	}
}
