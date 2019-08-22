package main

import (
	"AutoTrading/strategy"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":8000", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello world\n")

	strategy.RunPlatformDiffStrategy()
}
