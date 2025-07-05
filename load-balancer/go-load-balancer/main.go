package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	s := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			io.WriteString(w, "Hello from backend server")
		}),
	}
	log.Fatal(s.ListenAndServe())
}
