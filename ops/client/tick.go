package client

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"ops/pb"
	"time"
)

func tickUpdateClientStat(c net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			logError("tickUpdateClientStatinner error!", err)
		}
	}()

	for {
		cpu, err := Cpu()
		if err != nil {
			cpu = err.Error()
		}

		mem, err := Mem()
		if err != nil {
			mem = err.Error()
		}

		disk, err := Disk()
		if err != nil {
			disk = err.Error()
		}

		machineStatProto := &pb.MachineStatProto{
			Cpu:  proto.String(cpu),
			Mem:  proto.String(mem),
			Disk: proto.String(disk),
		}

		obj := &pb.UpdateClientStatProto{
			MachineStatProto: machineStatProto,
		}

		data, err := proto.Marshal(obj)
		fmt.Println(data)
		if err != nil {
			logError("tickUpdateClientStat proto.Marshalerror!", err)
		} else {
			logInfo("tickUpdateWorkerStat")
			err = writeBytes(c, assemble2(data, JOB_ID_UPDATE_CLIENT_STAT))
			if err != nil {
				logError("tickUpdateClientStat WriteBytes时出错!", err)
				return
			}
		}
		time.Sleep(time.Second * 3)
	}
}

func tickResClear() {
	defer func() {
		if err := recover(); err != nil {
			logError("tickDelRes内部错误!", err)
		}
		go tickResClear()
	}()

	////////for {
	////////	server.ResClear()
	////////	time.Sleep(time.Minute * 60)
	////////	// time.Sleep(time.Second * 5)
	////////}
}
