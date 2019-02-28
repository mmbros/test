package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	address := ":8082"

	backendURL := "http://127.0.0.1:8081"

	// parse th backend URL\
	u, err := url.Parse(backendURL)
	if err != nil {
		log.Fatal(err)
		//os.exit(1)
	}

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(u)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Update the headers to allow for SSL redirection
		r.URL.Host = u.Host
		r.URL.Scheme = u.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = u.Host

		// Note that ServeHttp is non blocking and uses a go routine under the hood
		proxy.ServeHTTP(w, r)
	})
	log.Printf("Listen and serve on %s", address)

	log.Fatal(http.ListenAndServe(address, logRequest(http.DefaultServeMux)))
}
