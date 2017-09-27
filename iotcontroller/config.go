package main

import (
	"encoding/json"
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
	Level  string `json:"level"`
	Url    string `json:"url"`
	Cron   string `json:"cron"`
	Sleep  int    `json:"sleep,string"`
	Apikey string `json:"apikey"`
	Limits []Limit
}

type Limit struct {
	Lower string `json:"lower"`
	Upper string `json:"upper"`
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
