package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// define helloWorld handler
func helloWorld(w http.ResponseWriter, req *http.Request) {
	hostname, err := os.Hostname()
	msg := fmt.Sprintf("Hello World from %s", hostname)
	if err != nil {
		io.WriteString(w, msg)
	}
	io.WriteString(w, msg)
}

// define main function
func main() {
	// define vars
	httpPort := ":8080"

	// initialise new servemux and register http handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloWorld)

	// start server
	log.Printf("Starting server on %s", httpPort)
	err := http.ListenAndServe(httpPort, mux)
	log.Fatal(err)
}
