package mokylin_client

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"mokylin_common/msg"
	"mokylin_common/msg/pb"
	"net"
	"time"
)

func tickUpdateClientStat(c net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("tick update clientstast error!", err)
		}
	}()
	for {
		cid := &pb.MachineStatProto{
			Cpu: proto.String("20%"),
		}
		obj := &pb.UpdateClientStatProto{
			MachineStatProto: cid,
		}
		data, err := proto.Marshal(obj)
		if err != nil {
			fmt.Println("proto.Marshal error!", err)
		} else {
			err := msg.WriteBytes(c, msg.Assemble2(data, msg.JOB_ID_UPDATE_CLIENT_STAT))
			if err != nil {
				fmt.Println("tick update clientstat msg.WriteBytes error!", err)
				return
			}
		}
		// tickup time is import, do not set too long, should less then 30s
		time.Sleep(time.Second * 30)
	}
}
