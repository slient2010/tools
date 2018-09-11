package mokylin_client

import (
	"flag"
	"fmt"
	//"log"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"strconv"
)

// serverinfos, databasesinfos come from config.ini keys
type Cfg struct {
	Server      []Server     `json:"server"`
	ClientInfos []ClientInfo `json:"clientinfos"`
}

// ipadds, port is the same as config.ini keys, be care with type
type Server struct {
	IpAddress  string `json:"ipadds"`
	ServerPort int    `json:"port"`
}

// dbserver, dbport, dbpass is the same as config.ini keys, be care with type
type ClientInfo struct {
	Id       string `json:"id"`
	Domain   string `json:"domain"`
	Gamepath string `json:"gamepath"`
}

// global variables
var (
	flags = flag.NewFlagSet("Client", flag.ExitOnError)
	cfg   string

	// service ip and port
	serverIpAndPort string

	// database info
	clientId string
	domain   string
	gamepath string
	ips      []string
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
	for _, si := range c.Server {
		serverIpAndPort = si.IpAddress + ":" + strconv.Itoa(si.ServerPort)
	}
	for _, ci := range c.ClientInfos {
		clientId = ci.Id
		domain = ci.Domain
		gamepath = ci.Gamepath
	}
}

// func externalIP() ([]string, error) {
func externalIP() (map[string]string, error) {
	ips := make(map[string]string)
	iplists := []string{}
	ifaces, err := net.Interfaces()
	if err != nil {
		//return iplists, err
		return ips, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return ips, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			iplists = append(iplists, ip.String())
		}
	}
	if len(iplists) > 0 {
		// get local ip and send to server
		lang, err := json.Marshal(iplists)
		if err != nil {
			fmt.Println(err)
		}
		//err = msg.WriteString(c, string(lang))
		if err != nil {
			fmt.Println(err)
		}
		//		ips["ip"] = string(lang)
		ips = map[string]string{"ip": string(lang)}
		//return iplists, nil
		return ips, err
	} else {
		// return iplists, errors.New("are you connected to the network?")
		return ips, errors.New("are you connected to the network?")
	}
}
