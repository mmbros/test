package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/mmbros/mananno/transmission"
)

var (
	trans *transmission.Client
)

func handlerTransmission(w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse(trans.Address)
	u.Path = r.URL.Path
	s := u.String()

	resp, err := trans.Get(s)
	if err != nil {
		log.Printf("!!! Get error !!!\n")
		log.Fatal(err)
	}

	//	if strings.HasSuffix(s, ".css") {
	//		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	//	}
	log.Printf("Content-Type = %s", w.Header().Get("Content-Type"))

	defer resp.Body.Close()
	text, err := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, "\n\n%s\n\n", text)
	if err != nil {
		log.Fatal(err)
	}
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf(">>>>>>>>> %s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	address := ":8082"

	trans = transmission.NewClient(
		"192.168.1.2:9091/transmission/",
		"",
		"")

	// create the reverse proxy
	//	proxy := httputil.NewSingleHostReverseProxy(u)

	// server1: DON'T REWRITE BODY
	http.HandleFunc("/transmission/", handlerTransmission)

	log.Printf("Listen and serve on %s", address)

	log.Fatal(http.ListenAndServe(address, logRequest(http.DefaultServeMux)))
}
