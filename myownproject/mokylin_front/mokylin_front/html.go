package mokylin_front

import (
	"mokylin_common/msg/pb"
)

type HtmlClient struct {
	Id          string
	Path        string
	ServerCount int
	ServerNames string
	Cpu         string
	Mem         string
	Disk        string
	Stat        string
}

type CreateServerJob struct {
	clientId  string
	jobProto  *pb.CreateServerJobProto
	replyChan chan []byte
}
