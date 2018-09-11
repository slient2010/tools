package master

// peer提交给manager的
type clientReg struct {
	client *client
	reply  chan bool
}

// manager提交给peer的
type job struct {
	data      []byte
	replyChan chan []byte
}
