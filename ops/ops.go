package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"ops/client"
	"os"
	"strconv"
	"strings"
)

const (
	fileCfg = "config.yml"
)

type Cfg struct {
	Master struct {
		Ip   string `yaml:"ip"`
		Port int    `yaml:"port"`
	}
}

func main() {
	flag.Parse()
	action := flag.Arg(0)
	if strings.HasPrefix(action, "start") {
		initCfg()
		client.Start()
	} else if action == "stop" {
		client.Stop()
	} else {
		usage()
	}

	// fmt.Println("vim-go")
}

func usage() {
	fmt.Println("ops client usage!")
}

func initCfg() {
	data, err := ioutil.ReadFile(fileCfg)
	if err != nil {
		if os.IsNotExist(err) {
			os.Exit(1)
			return
		}
		log.Fatal(err)
	}
	var cfg Cfg
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	client.MasterAddr = cfg.Master.Ip + ":" + strconv.Itoa(cfg.Master.Port)
	fmt.Println(cfg)
}
