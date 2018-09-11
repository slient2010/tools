package mokylin_server

import (
	"flag"
	"fmt"
	//"log"
	"encoding/json"
	"io/ioutil"
	"strconv"
)

// serverinfos, databasesinfos come from config.ini keys
type Cfg struct {
	ServerInfos   []ServerInfo   `json:"serverinfos"`
	DatabaseInfos []DatabaseInfo `json:"databaseinfos"`
}

// ipadds, port is the same as config.ini keys, be care with type
type ServerInfo struct {
	IpAddress  string `json:"ipadds"`
	ServerPort int    `json:"port"`
	HttpListen string `json:"httplisten"`
	HttpPort   int    `json:"httpport"`
}

// dbserver, dbport, dbpass is the same as config.ini keys, be care with type
type DatabaseInfo struct {
	DbServer string `json:"dbserver"`
	DbPort   int    `json:"dbport"`
	DbUser   string `json:"dbuser"`
	DbPass   string `json:"dbpass"`
}

// global variables
var (
	flags = flag.NewFlagSet("Server", flag.ExitOnError)
	cfg   string

	// service ip and port
	listenIpAndPort string
	// http service ip and port
	httpListenIpAndPort string

	// database info
	databaseip   string
	databaseport int
	databaseuser string
	databasepass string
)

func init() {
	// load configuration file
	flags.StringVar(&cfg, "c", "config/config.ini", "server config")

	data, err := ioutil.ReadFile(cfg)
	if err != nil {
		//    return err
		fmt.Println(err)
	}

	var c Cfg
	err = json.Unmarshal(data, &c)
	if err != nil {
		//		return err
		fmt.Println(err)
	}
	for _, si := range c.ServerInfos {
		listenIpAndPort = si.IpAddress + ":" + strconv.Itoa(si.ServerPort)
		httpListenIpAndPort = si.HttpListen + ":" + strconv.Itoa(si.HttpPort)
	}
	for _, db := range c.DatabaseInfos {
		databaseip = db.DbServer
		databaseport = db.DbPort
		databaseuser = db.DbUser
		databasepass = db.DbPass
	}
}

func Start() {
	// get data and show on the web
	go GetDatabaseInfo()
	// start services
	clientManager(listenIpAndPort)
	// start http services
//	front.HttpServer(httpListenIpAndPort)
}
