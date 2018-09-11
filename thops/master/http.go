package master

import (
	"log"
	"net/http"
)

func httpServer(addr string) {
	http.HandleFunc("/ping", ping)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
