package main

import (
	"io"
	"log"
	"net/http"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello world\n")
}

func main() {
	http.HandleFunc("/hello", HelloServer)

	err := http.ListenAndServe(":8888", nil)

	if err != nil {
		log.Fatal("ListenAndServer: ", err.Error())
	}
}
