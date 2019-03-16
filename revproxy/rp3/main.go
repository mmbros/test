package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

const headerTransmissionSessionID = "X-Transmission-Session-Id"

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

type transport struct {
	http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	log.Printf(resp.Status)
	if resp.StatusCode == http.StatusConflict {
		transmissionSessionID := resp.Header.Get(headerTransmissionSessionID)
		log.Printf("%s = %s", headerTransmissionSessionID, transmissionSessionID)
		req.Header.Add(headerTransmissionSessionID, transmissionSessionID)

		resp, err = t.RoundTripper.RoundTrip(req)
		if err != nil {
			return nil, err
		}
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1)
	body := ioutil.NopCloser(bytes.NewReader(b))

	resp.Body = body
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
	return resp, nil
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	address := ":8082"

	backendURL := "http://192.168.1.2:9091"

	// parse th backend URL\
	u, err := url.Parse(backendURL)
	if err != nil {
		log.Fatal(err)
		//os.exit(1)
	}

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = &transport{http.DefaultTransport}

	// server1: DON'T REWRITE BODY
	http.HandleFunc("/transmission/", func(w http.ResponseWriter, r *http.Request) {
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
