package mokylin_server

import (
	"fmt"
	"mokylin_common/msg/pb"
)

const (
	chanBufSize = 1000
)

var (
	// register channel
	clientRegChan = make(chan *clientReg, chanBufSize)
	// unregister channel
	clientUnRegChan = make(chan string, chanBufSize)
	//update the client status
	updateClientStatChan = make(chan *pb.UpdateClientStatProto, chanBufSize)
)

func clientManager(listenIpAndPort string) {
	go clientListener(listenIpAndPort)
	fmt.Println(listenIpAndPort)

	clientMap := make(map[string]*client)

	for {
		select {
		case reg := <-clientRegChan:
			processClientReg(clientMap, reg)
		// unregister client
		case clientId := <-clientUnRegChan:
			fmt.Println("--unregister client", clientId)
			processClientUnReg(clientMap, clientId)
		case o := <-updateClientStatChan:
			processUpdateClientStat(clientMap, o)
		}
	}

}

func registerClient(w *client) bool {
	okChan := make(chan bool, 1)
	clientRegChan <- &clientReg{w, okChan}
	return <-okChan
}

func unregisterClient(w *client) {
	clientUnRegChan <- w.Id
}

// sent update command into channel
func updateClientStat(o *pb.UpdateClientStatProto) {
	updateClientStatChan <- o
}

func processClientReg(clientMap map[string]*client, wr *clientReg) {
	client := wr.client
	if _, ok := clientMap[client.Id]; ok {
		fmt.Println("client已经存在，但又收到了注册申请!", client)
		wr.reply <- false
		return
	}
	clientMap[client.Id] = client
	wr.reply <- true
	// logInfo(reflect.TypeOf(client))
	// save client into db
	// onClientConnected(client.Id)
	fmt.Println("save client into db")
	////////checkres, err := CheckClient(client.Id, client.Domain)
	////////if err != nil {
	////////	// logInfo(err)
	////////	fmt.Println(err)
	////////}
	////////if checkres == 0 {
	fmt.Println("save client")
	err := SaveClientInfo(client.Id, client.Ipaddress, client.Domain, client.Path)
	if err == nil {
		fmt.Println("save client done")
	} else {
		fmt.Println(err)
	}
	////////	}
	fmt.Println("client注册成功, 需完成入库操作...", client.Id)
}

// process client unregister
func processClientUnReg(clientMap map[string]*client, clientId string) {
	// if client, ok := clientMap[clientId]; !ok {
	client, _ := CheckClient(clientId, clientMap[clientId].Domain)
	//	if client, ok := CheckClient(clientId, clientMap[clientId].Domain); !ok {
	if client != 1 {
		fmt.Println("client cancel register, but not found!", clientId)
	} else {
		UpdateClientUnRegister(clientId, clientMap[clientId].Domain)
		//delete the clientId, will not use this method later
		c := clientMap[clientId]
		delete(clientMap, clientId)
		//close the channel
		close(c.waitingSendJobsChan)
		fmt.Println("client unregistered!", clientId)
	}
}

// update client state
func processUpdateClientStat(clientMap map[string]*client, o *pb.UpdateClientStatProto) {
	client, _ := clientMap[*o.ClientId]
	err := SaveClientInfo(client.Id, client.Ipaddress, client.Domain, client.Path)
	//	_, ok := clientMap[*o.ClientId]
	if err != nil {
		fmt.Println("update client, but not found!", o.ClientId)
		return
	}
}
