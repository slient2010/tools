package mokylin_server

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"mokylin_common/msg"
	"mokylin_common/msg/pb"
	"net"
	"time"
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

func clientListener(addr string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("client listener accept error!", err)
			continue
		}
		go clientHandler(conn)
	}
}

func clientHandler(c net.Conn) {
	// get the clientid
	clientId, err := msg.ReadLimitString(c, 100)
	if err != nil {
		fmt.Println("read client error!", err)
		c.Close()
		return
	}

	// get client ipaddress
	ips, err := msg.ReadString(c)
	if err != nil {
		fmt.Println("read client ipaddress error!", err)
		c.Close()
		return
	}

	// get the client domain
	domain, err := msg.ReadString(c)
	if err != nil {
		fmt.Println("read client domain error!", err)
		c.Close()
		return
	}

	// get the gamepath
	gamePath, err := msg.ReadString(c)
	if err != nil {
		fmt.Println("read client game path error!", err)
		c.Close()
		return
	}

	//
	waitingSendJobsChan := make(chan *job, 100)

	client := &client{
		Id:                  clientId,
		Ipaddress:           ips,
		Domain:              domain,
		Path:                gamePath,
		waitingSendJobsChan: waitingSendJobsChan,
		conn:                c,
	}

	if !registerClient(client) {
		msg.WriteBool(c, false)
		c.Close()
		return
	}

	//
	err = msg.WriteBool(c, true)
	if err != nil {
		fmt.Println(err)
		c.Close()
		return
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("inner error!", err)
			// delete register
			unregisterClient(client)
			c.Close()
		}
	}()

	jobIdCounter := 0
	waitingReplyJobsMap := make(map[int]chan []byte)
	jobsReplyChan := make(chan []byte, 10)
	clientCloseChan := make(chan int)
	heartBeatTimeOut := false

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("clientHandler读数据内部错误！", err)
				clientCloseChan <- 1
			}
		}()

		for {
			// recive tick up info msg
			data, err := msg.ReadBytes(c)
			if err != nil {
				clientCloseChan <- 1
				return
			}
			jobsReplyChan <- data

			heartBeatTimeOut = false
		}
	}()

	// keep connecting with clients, otherwise close the connections.
	go func() {
		tickChan := time.Tick(time.Second * 30)
		for _ = range tickChan {
			if heartBeatTimeOut {
				fmt.Println("clientHandler heartBeatTimeOut:", clientId)
				clientCloseChan <- 1
			} else {
				heartBeatTimeOut = true
			}
		}
	}()

	// communicate with client
	for {
		select {
		case job := <-waitingSendJobsChan:
			jobIdCounter += 1
			jobId := jobIdCounter & msg.JOB_ID_MASK
			if _, ok := waitingReplyJobsMap[jobId]; ok {
				fmt.Println("clienthandler get jobId error!", jobId)
				closeClient(c, client, waitingReplyJobsMap, jobsReplyChan)
				return
			}
			jobData := msg.Assemble2(job.data, jobId)

			err := msg.WriteBytes(c, jobData)
			if err != nil {
				fmt.Println("clientHandler写client job时出错！", err)
				closeClient(c, client, waitingReplyJobsMap, jobsReplyChan)
				return
			}

			waitingReplyJobsMap[jobId] = job.replyChan
		case replyData := <-jobsReplyChan:
			data, jobId := msg.Split2(replyData)

			switch jobId {
			case msg.JOB_ID_UPDATE_CLIENT_STAT:
				fmt.Println("Update client stat")
				obj := &pb.UpdateClientStatProto{}
				err := proto.Unmarshal(data, obj)
				if err != nil {
					fmt.Println("收到更新服务器状态的消息，但unmarshal proto时出错!", err)
					continue
				}
				obj.ClientId = proto.String(clientId)
				// update status
				updateClientStat(obj)
				//			case msg.JOB_ID_UNION_SERVER:
				//				submitUnionServerJobReply(msg.AssembleUtfStr(data, clientId))
			default:
				if replyChan, ok := waitingReplyJobsMap[jobId]; ok {
					replyChan <- data
					delete(waitingReplyJobsMap, jobId)
				} else {
					fmt.Println("收到jobReply，但在waitingReplyJobsMap中没有找到对应的job!", jobId, len(data))
				}
			}

		case <-clientCloseChan:
			closeClient(c, client, waitingReplyJobsMap, jobsReplyChan)
			return
		}
	}

}

func GetDatabaseInfo() {
	// get client infos and show them on web
	fmt.Println("database", databaseip)
}

func closeClient(c net.Conn, w *client, waitingReplyJobsMap map[int]chan []byte, jobsReplyChan chan []byte) {
	c.Close()
	unregisterClient(w)

	close(jobsReplyChan)
	for replyData := range jobsReplyChan {
		data, jobId := msg.Split2(replyData)
		if replyChan, ok := waitingReplyJobsMap[jobId]; ok {
			replyChan <- data
			delete(waitingReplyJobsMap, jobId)
		}
	}

	clientClosed := []byte("client连接中断!")
	for _, replyChan := range waitingReplyJobsMap {
		replyChan <- clientClosed
	}

	for job := range w.waitingSendJobsChan { // waitingSendJobsChan closed by unregister
		job.replyChan <- clientClosed
	}
}
