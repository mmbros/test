// Golang: Log HTTP Requests in Go
// https://gist.github.com/hoitomt/c0663af8c9443f2a8294

package main

import (
	"fmt"
	"log"
	"net/http"
)

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	address := ":8081" // don't change (otherwise update index.html)

	http.Handle("/", http.FileServer(http.Dir("./static")))
	//
	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi")

	})

	log.Printf("Listen and serve on %s", address)

	log.Fatal(http.ListenAndServe(address, logRequest(http.DefaultServeMux)))
}
