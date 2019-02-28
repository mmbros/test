// https://golang.org/pkg/net/http/httputil/#example_ReverseProxy

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
)

func main() {
	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "this call was relayed by the reverse proxy")

	}))
	defer backendServer.Close()

	log.Printf("backendServer.UTL = %s\n", backendServer.URL)

	rpURL, err := url.Parse(backendServer.URL)
	if err != nil {
		log.Fatal(err)

	}
	frontendProxy := httptest.NewServer(httputil.NewSingleHostReverseProxy(rpURL))
	defer frontendProxy.Close()

	resp, err := http.Get(frontendProxy.URL)
	if err != nil {
		log.Fatal(err)

	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)

	}

	fmt.Printf("%s", b)

}
