package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

// https://stackoverflow.com/questions/31535569/golang-how-to-read-response-body-of-reverseproxy
// httputil.ReverseProxy has a Transport field. You can use it to modify the response. For example:

type transport struct {
	http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	b = bytes.Replace(b, []byte("http://localhost:8081/"), []byte(""), -1)
	// make path from absolute to relative
	b = bytes.Replace(b, []byte("=\"/"), []byte("=\"./"), -1)

	body := ioutil.NopCloser(bytes.NewReader(b))
	resp.Body = body
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
	return resp, nil
}

func main() {
	address := ":8082"

	//backendURL := "http://127.0.0.1:8081"
	backendURL := "http://192.168.1.2:8081"
	//backendURL := "http://mananno.dlinkddns.com:9091"

	// parse th backend URL\
	u, err := url.Parse(backendURL)
	if err != nil {
		log.Fatal(err)
		//os.exit(1)
	}

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = &transport{http.DefaultTransport} // for server2

	// server1: DON'T REWRITE BODY
	http.HandleFunc("/server1/", func(w http.ResponseWriter, r *http.Request) {
		// Update the headers to allow for SSL redirection
		r.URL.Host = u.Host
		r.URL.Scheme = u.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = u.Host
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/server1")

		// Note that ServeHttp is non blocking and uses a go routine under the hood
		proxy.ServeHTTP(w, r)
	})

	// server2: REWRITE BODY
	http.HandleFunc("/server2/", func(w http.ResponseWriter, r *http.Request) {

		// Update the headers to allow for SSL redirection
		r.URL.Host = u.Host
		r.URL.Scheme = u.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = u.Host
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/server2")

		// Note that ServeHttp is non blocking and uses a go routine under the hood
		proxy.ServeHTTP(w, r)
	})

	log.Printf("Listen and serve on %s", address)

	log.Fatal(http.ListenAndServe(address, logRequest(http.DefaultServeMux)))
}
