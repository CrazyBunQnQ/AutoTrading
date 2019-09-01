package main

import (
	"AutoTrading/strategy"
	"io"
	"net/http"
)

func main() {
	//http.HandleFunc("/", start)
	//http.HandleFunc("/", stop)
	//http.ListenAndServe(":8000", nil)
	platformDiff()
}

func start(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "start\n")
}

func stop(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "stop\n")
}

func platformDiff() {
	strategy.RunPlatformDiffStrategy()
}
