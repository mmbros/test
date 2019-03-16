package main

import (
	// "fmt"
	// "io"
	"fmt"
	"log"
	"net/http"
)

// HelloServer is ...
func HelloServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
	// fmt.Fprintf(w, "This is an example server.\n")
	// io.WriteString(w, "This is an example server.\n")

}

func main() {
	port := ":4443"
	http.HandleFunc("/", HelloServer)
	fmt.Printf("Listening on %s ...\n", port)
	err := http.ListenAndServeTLS(port, "mananno.cert.pem", "mananno.key.pem", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
