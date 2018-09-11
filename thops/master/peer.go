package master

import (
	"fmt"
	"log"
	"net"
	"thops/pb"
	"time"
)

type client struct {
	Id                  string
	Path                string               // work path
	MachineStatProto    *pb.MachineStatProto // Machine property
	waitingSendJobsChan chan *job
	conn                net.Conn
}

func clientListener(addr, secret string) {
	// fmt.Println("master listener start!", addr)
	logInfo("master listener start!", addr)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
		// logError(err)
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			logError("client listener accept error!", err)
			continue
		}
		go clientHandler(conn)
	}
}

func clientHandler(c net.Conn) {
	defer func() {
		fmt.Println(c.RemoteAddr(), "Closed")
		c.Close()
	}()

	clientId, err := readString(c)
	if err != nil {
		logError("read clientid error!", err)
		c.Close()
		return
	}
	workPath, err := readString(c)
	if err != nil {
		logError("read client work path error!", err)
		c.Close()
		return
	}

	// client struct
	client := &client{
		Id:               clientId,
		Path:             workPath,
		MachineStatProto: &pb.MachineStatProto{},
		conn:             c,
	}

	if !registerClient(client) {
		writeBool(c, false)
		c.Close()
		return
	}

	err = writeBool(c, true)
	if err != nil {
		logError("save client error!", err)
		c.Close()
		return
	}

	// jobIdCounter := 0
	waitingReplyJobsMap := make(map[int]chan []byte)

	jobsReplyChan := make(chan []byte, 10)
	workerCloseChan := make(chan int)

	heartBeatTimeOut := false

	go func() {
		defer func() {
			if err := recover(); err != nil {
				logError("workerHandler读数据内部错误！", err)
				workerCloseChan <- 1
			}
		}()

		for {
			data, err := readBytes(c)
			if err != nil {
				workerCloseChan <- 1
				return
			}
			jobsReplyChan <- data

			heartBeatTimeOut = false
		}
	}()

	go func() {
		tickChan := time.Tick(time.Second * 30)
		for _ = range tickChan {
			if heartBeatTimeOut {
				logError("workerHandler heartBeatTimeOut:", clientId)
				workerCloseChan <- 1
			} else {
				heartBeatTimeOut = true
			}
		}
	}()

	for {
		select {
		case replyData := <-jobsReplyChan:
			fmt.Println(replyData)
			// data, jobId := Split2(replyData)

		////////switch jobId {
		////////case msg.JOB_ID_UPDATE_WORKER_STAT:
		////////	obj := &pb.UpdateWorkerStatProto{}
		////////	err := proto.Unmarshal(data, obj)
		////////	if err != nil {
		////////		logError("收到更新服务器状态的消息，但unmarshal proto时出错!", err)
		////////		continue
		////////	}
		////////	obj.WorkerId = proto.String(workerId)
		////////	updateWorkerStat(obj)
		////////case msg.JOB_ID_UNION_SERVER:
		////////	submitUnionServerJobReply(msg.AssembleUtfStr(data, workerId))
		////////default:
		////////	if replyChan, ok := waitingReplyJobsMap[jobId]; ok {
		////////		replyChan <- data
		////////		delete(waitingReplyJobsMap, jobId)
		////////	} else {
		////////		logError("收到jobReply，但在waitingReplyJobsMap中没有找到对应的job!", jobId, len(data))
		////////	}
		////////}

		case <-workerCloseChan:
			closeWorker(c, client, waitingReplyJobsMap, jobsReplyChan)
			return
		}
	}

}

func closeWorker(c net.Conn, client *client, waitingReplyJobsMap map[int]chan []byte, jobsReplyChan chan []byte) {
	c.Close()
	// unregisterWorker(w)

	close(jobsReplyChan)
	////////for replyData := range jobsReplyChan {
	////////	data, jobId := msg.Split2(replyData)
	////////	if replyChan, ok := waitingReplyJobsMap[jobId]; ok {
	////////		replyChan <- data
	////////		delete(waitingReplyJobsMap, jobId)
	////////	}
	////////}

	clientClosed := []byte("worker连接中断!")
	for _, replyChan := range waitingReplyJobsMap {
		replyChan <- clientClosed
	}

	for job := range client.waitingSendJobsChan { // waitingSendJobsChan closed by unregister
		job.replyChan <- clientClosed
	}
}
