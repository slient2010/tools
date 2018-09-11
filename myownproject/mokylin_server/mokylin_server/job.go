package mokylin_server

type clientReg struct {
	client *client
	reply  chan bool
}

type job struct {
	data      []byte
	replyChan chan []byte
}
