package main

import (
	"net/http"
)

func root(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("hello\n"))
}

func main() {
	http.HandleFunc("/", root)
	http.ListenAndServe(":8080", nil)
}
