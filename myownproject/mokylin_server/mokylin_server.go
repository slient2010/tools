package main

import (
    "mokylin_server/mokylin_server"
    "fmt"
    "runtime"
)

func main() {
    fmt.Println("Start mokylin server")
    runtime.GOMAXPROCS(runtime.NumCPU())
    mokylin_server.Start()
}
