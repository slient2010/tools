package client

import (
	"fmt"
	"log"
	"net"
	"time"
)

var (
	conn       net.Conn
	MasterAddr string
)

func Start() {
	fmt.Println("start client")
	conn = connectUntilSuccess()
	go tickUpdateClientStat(conn)
	for {
		data, err := readBytes(conn)
		fmt.Println(data)
		if err != nil {
			conn.Close()
		}
		go process(conn, data)
	}
}

func Stop() {
	fmt.Println("stop client")
}

func connectUntilSuccess() net.Conn {
	for {
		c, err := connect()
		if err == nil {
			return c
		}
		time.Sleep(time.Second)
	}
}

func connect() (net.Conn, error) {
	logInfo("try connect master:", MasterAddr)

	c, err := net.Dial("tcp", MasterAddr)
	if err != nil {
		return nil, err
	}

	clientId := "client-172.19.1.1"
	err = writeString(c, clientId)
	if err != nil {
		return nil, err
	}

	baseDir := "/data/tehang/apps"
	err = writeString(c, baseDir)
	if err != nil {
		return nil, err
	}

	ok, err := readBool(c)
	if err != nil {
		return nil, err
	}

	if !ok {
		log.Fatal("client注册失败!")
	}

	logInfo("client连接成功!", clientId)

	return c, nil
}

func process(c net.Conn, data []byte) {
	fmt.Println(data)
	// jobInfo, jobId := split2(data)

	jobId := 1
	jobInfo := "aa"
	logInfo("process job:", jobId)

	fmt.Println(jobId)
	fmt.Println(jobInfo)
	////////switch jobId {
	////////case JOB_ID_UNION_SERVER:
	////////	processUnionServer(c, jobInfo)
	////////default:
	////////	jobData, jobType := msg.Split(jobInfo)

	////////	switch jobType {
	////////	case msg.JOB_TYPE_CREATE_SERVER:
	////////		processCreateServer(c, jobId, jobData)
	////////	case msg.JOB_TYPE_START_SERVER:
	////////		processStartServer(c, jobId, jobData)
	////////	case msg.JOB_TYPE_STOP_SERVER:
	////////		processStopServer(c, jobId, jobData)
	////////	case msg.JOB_TYPE_UPDATE_SERVER:
	////////		processUpdateServer(c, jobId, jobData)
	////////	default:
	////////		msg.WriteBytes(c, msg.Assemble2([]byte("illegal jobType:"+strconv.Itoa(jobType)), jobId))
	////////		logError("illegal jobType:", jobType)
	////////	}
	////////}
}
