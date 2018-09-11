package front

import (
	"flag"
	//    "fmt"
)

var (
	flags    = flag.NewFlagSet("front", flag.ExitOnError)
	portHttp string
)

func init() {
	flags.StringVar(&portHttp, "ph", "0.0.0.0:8088", "front http service port")
}

// func Start(args []string) {
//    flag.Parse(args)
func Start() {
	httpServer(portHttp)
}
