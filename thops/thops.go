package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"thops/master"
	_ "thops/server"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	action := flag.Arg(0)
	if strings.HasPrefix(action, "start") {
		fmt.Println("Loading configurations!")
		initCfg()
		master.Start()
	} else if action == "stop" {
		master.Stop()
	} else {
		Usage()
	}
}

func Usage() {
	fmt.Println("devops usage!")
}

const (
	fileCfg = "config.yml"
)

type Cfg struct {
	Server struct {
		Address string `yaml:"address"`
		WebPort int    `yaml:"webport"`
		Port    int    `yaml:"port"`
		Secret  string `yaml:"secret"`
	}
	Database struct {
		Ip       string `yaml:"ip"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Dbname   string `yaml:"dbname"`
		Charset  string `yaml:"charset"`
	}
	RedisDb struct {
		Ip       string `yaml:"ip"`
		Port     int    `yaml:"port"`
		Password string `yaml:"password"`
	}
	Resource struct {
		Ip   string `yaml:"ip"`
		Port int    `yaml:"port"`
	}
}

func initCfg() {
	data, err := ioutil.ReadFile(fileCfg)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Configuration not exist!")
			return
		}
		log.Fatal(err)
	}
	var cfg Cfg
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		fmt.Println("configuration wrong!")
		log.Fatal(err)
	}

	master.PortClientListener = cfg.Server.Port
	master.MasterAddress = cfg.Server.Address
	master.WebPort = cfg.Server.WebPort
	master.Secret = cfg.Server.Secret

	master.Host = cfg.Database.Ip
	master.Port = cfg.Database.Port
	master.Username = cfg.Database.Username
	master.Password = cfg.Database.Password
	master.Dbname = cfg.Database.Dbname
	master.Charset = cfg.Database.Charset

	// fmt.Println(cfg)
}
