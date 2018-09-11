package mokylin_front

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	_ "mokylin_common/msg"
	"mokylin_common/msg/pb"
	"net"
)

const (
	chanBufSize = 1000
)

type client struct {
	Id string
	// game ip address
	Ipaddress string
	// client domain
	Domain string
	// game path
	Path string
	// machine status pointer
	// MachineStatProto *pb.MachineStatProto
	// ServerStatProtoMap map[string]*pb.ServerStatProto
	waitingSendJobsChan chan *job
	conn                net.Conn
}

type job struct {
	data      []byte
	replyChan chan []byte
}

var (
	getClientsChan      = make(chan chan *client, chanBufSize)
	createServerJobChan = make(chan *CreateServerJob, chanBufSize)
)

func JobManager() {
	fmt.Println("job listener")
	clientMap := make(map[string]*client)
	for {
		select {
		case o := <-getClientsChan:
			processGetClients(clientMap, o)
		case o := <-createServerJobChan:
			processCreateServer(clientMap, o)

		}
	}
}

func processCreateServer(clientMap map[string]*client, o *CreateServerJob) {
	clientId := o.clientId
	jobProto := o.jobProto
	replyChan := o.replyChan

	////////client, ok := clientMap[clientId]
	////////if ok {
	////////	replyChan <- []byte("clientId不存在!" + clientId)
	////////	return
	////////}

	// send
	data, err := proto.Marshal(jobProto)
	if err != nil {
		replyChan <- []byte("内部错误!")
	} else {
		fmt.Println(clientId)
		fmt.Println(jobProto)
		fmt.Println(data)
		return
                // use api send job
	////////client.waitingSendJobsChan <- &job{
	////////	msg.Assemble(data, msg.JOB_TYPE_CREATE_SERVER),
	////////	replyChan,
	////////}
	}
}

func submitCreateServerJob(clientId string, jobProto *pb.CreateServerJobProto) []byte {
	replyChan := make(chan []byte, 1)
	job := &CreateServerJob{
		clientId,
		jobProto,
		replyChan,
	}
	createServerJobChan <- job
	return <-replyChan
}

func processGetClients(clientMap map[string]*client, reqChan chan *client) {
	for _, client := range clientMap {
		reqChan <- client
	}
	close(reqChan)
}
