package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/testEndPoint", customHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error occured unable to started: %v", err)
	}
}

func customHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello from backend server")
}
