package main

import (
	"fmt"
	"mokylin_front/mokylin_front"
	"runtime"
)

func main() {
	fmt.Println("start http server")
	runtime.GOMAXPROCS(runtime.NumCPU())
	go mokylin_front.JobManager()
	mokylin_front.HttpServer("0.0.0.0:80")
}
