package main

import (
	"AutoTrading/strategy"
	"io"
	"net/http"
	"time"
)

func main() {
	//http.HandleFunc("/", start)
	//http.HandleFunc("/", stop)
	//http.ListenAndServe(":8000", nil)
	for true {
		time.Sleep(time.Duration(5) * time.Second)
		strategy.RunPlatformDiffStrategy()
	}
}

func start(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "start\n")
}

func stop(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "stop\n")
}
