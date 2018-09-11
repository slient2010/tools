package master

import (
	"fmt"
)

const (
	chanBufSize = 1000
)

var (
	clientRegChan = make(chan *clientReg, chanBufSize)
)

func clientManager(listenAddr, secret string) {
	go clientListener(listenAddr, secret)

	for {
		select {
		case reg := <-clientRegChan:
			processClientReg(reg)
		}
	}

}

func registerClient(client *client) bool {
	okChan := make(chan bool, 1)
	clientRegChan <- &clientReg{client, okChan}
	return <-okChan
}

func processClientReg(cr *clientReg) {
	client := cr.client
	data := &ServerInfo{Name: "tehang-dev-server", Ipaddress: "172.19.1.11", CPU: 12, Mem: 8, Disk: 500, UUID: "test"}
	// check client
	ok := checkClient(data)
	if ok {
		logError("client already registered!", client)
		cr.reply <- false
		return
	}
	saveClient(client, data)
	cr.reply <- true
	onClientConnected(client.Id)
	logError("Client register success!", client.Id)
}

func checkClient(d *ServerInfo) bool {
	count := checkServerInOrNotInDb(d)
	if count >= 1 {
		return true
	}
	return false
}

func saveClient(p *client, data *ServerInfo) {
	insertIntoDb(data)
	fmt.Println(p)
}

func onClientConnected(clientId string) {
	fmt.Println(clientId)
}
