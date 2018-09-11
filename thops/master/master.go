package master

import (
	"flag"
	"fmt"
	"strconv"
)

var (
	flags              = flag.NewFlagSet("master", flag.ExitOnError)
	PortClientListener int
	MasterAddress      string
	WebPort            int
	Secret             string
)

// func Start(args []string) {
func Start() {
	// flags.Parse(args)

	fmt.Println("master start")
	go clientManager(":"+strconv.Itoa(PortClientListener), Secret)
	httpServer(MasterAddress + ":" + strconv.Itoa(WebPort))
}

// func Stop(args []string) {
func Stop() {

	fmt.Println("master vim-go")
}
