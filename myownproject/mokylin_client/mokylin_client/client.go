package mokylin_client

import (
	"encoding/json"
	"fmt"
	"mokylin_common/msg"
	"net"
	"sync/atomic"
	"time"
        "strconv"
)

var (
	conn            net.Conn
	doingJobCounter int32
)

func Start() {
	//	fmt.Println(clientId)
	//	fmt.Println(domain)
	conn = connectUntilSuccess()
	go tickUpdateClientStat(conn)
	atomic.StoreInt32(&doingJobCounter, 0)
	for {
		data, err := msg.ReadBytes(conn)
		if err != nil {
			conn.Close()
			fmt.Println("client msg.ReadBytes error, reconnect!", err)
			conn = connectUntilSuccess()
			go tickUpdateClientStat(conn)
			continue
		}
		fmt.Println("recv job...")
		atomic.AddInt32(&doingJobCounter, 1)
		go process(conn, data)
	}

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
	fmt.Println("["+time.Now().Format("2006-01-02 15:04:05")+"] try connect master:", serverIpAndPort)
	c, err := net.Dial("tcp", serverIpAndPort)
	if err != nil {
		return nil, err
	}

	// clientid
	err = msg.WriteString(c, clientId)
	if err != nil {
		return nil, err
	}

	// get local ip and send to server
	ips, err := externalIP()
	if err != nil {
		fmt.Println(err)
	}

	lang, err := json.Marshal(ips)
	if err != nil {
		fmt.Println(err)
	}
	err = msg.WriteString(c, string(lang))
	if err != nil {
		return nil, err
	}

	// localhost domain
	err = msg.WriteString(c, domain)
	if err != nil {
		return nil, err
	}

	// game path
	err = msg.WriteString(c, gamepath)
	if err != nil {
		return nil, err
	}

	// if registered, then failed
	ok, err := msg.ReadBool(c)
	if err != nil {
		return nil, err
	}

	if !ok {
		fmt.Println("Client register failed!")
	}
	fmt.Println("["+time.Now().Format("2006-01-02 15:04:05")+"] connect master success!", clientId)
	return c, err

}

func process(c net.Conn, data []byte) {
	defer func() {
		atomic.AddInt32(&doingJobCounter, -1)
	}()

	jobInfo, jobId := msg.Split2(data)

	//logInfo("process job:", jobId)

	switch jobId {
////////case msg.JOB_ID_UNION_SERVER:
////////	processUnionServer(c, jobInfo)
	default:
		jobData, jobType := msg.Split(jobInfo)

               fmt.Println(jobData)
		switch jobType {
	////////case msg.JOB_TYPE_CREATE_SERVER:
	////////	processCreateServer(c, jobId, jobData)
	////////case msg.JOB_TYPE_START_SERVER:
	////////	processStartServer(c, jobId, jobData)
	////////case msg.JOB_TYPE_STOP_SERVER:
	////////	processStopServer(c, jobId, jobData)
	////////case msg.JOB_TYPE_UPDATE_SERVER:
	////////	processUpdateServer(c, jobId, jobData)
		default:
			msg.WriteBytes(c, msg.Assemble2([]byte("illegal jobType:"+strconv.Itoa(jobType)), jobId))
			fmt.Println("illegal jobType:", jobType)
		}
	}

	fmt.Println("=========")
	fmt.Println(c)
	fmt.Println("=========")
	fmt.Println(data)
}
