package main

// REFERENCES
// - https://github.com/jcbsmpsn/golang-https-example/blob/master/https_server.go
// - https://gist.github.com/xjdrew/97be3811966c8300b724deabc10e38e2

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

var (
	caFile = flag.String("CA", "ca.cert.pem", "A PEM eoncoded CA's certificate file.")

	certFile = flag.String("cert", "localhost.cert.pem", "A PEM eoncoded certificate file.")
	keyFile  = flag.String("key", "localhost.key.pem", "A PEM encoded private key file.")

	addr = ":4443"
)

func main01() {

	// redirect every http request to https
	go http.ListenAndServe(":8080", http.HandlerFunc(redirect))

	flag.Parse()

	caCert, err := ioutil.ReadFile(*caFile)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	cfg := &tls.Config{
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  caCertPool,
	}
	srv := &http.Server{
		Addr:           addr,
		Handler:        &handler{},
		TLSConfig:      cfg,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("listen on %s", addr)
	log.Fatal(srv.ListenAndServeTLS(*certFile, *keyFile))
	// log.Fatal(srv.ListenAndServe())
}

type handler struct{}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	printConnState(req.TLS)
	w.Write([]byte("PONG\n"))
}
func printConnState(state *tls.ConnectionState) {
	log.Print(">>>>>>>>>>>>>>>> State <<<<<<<<<<<<<<<<")
	if state == nil {
		log.Printf("STATE = NIL")
		return
	}
	log.Printf("Version: %x", state.Version)
	log.Printf("HandshakeComplete: %t", state.HandshakeComplete)
	log.Printf("DidResume: %t", state.DidResume)
	log.Printf("CipherSuite: %x", state.CipherSuite)
	log.Printf("NegotiatedProtocol: %s", state.NegotiatedProtocol)
	log.Printf("NegotiatedProtocolIsMutual: %t", state.NegotiatedProtocolIsMutual)

	log.Print("Certificate chain:")
	for i, cert := range state.PeerCertificates {
		subject := cert.Subject
		issuer := cert.Issuer
		log.Printf(" %d s:/C=%v/ST=%v/L=%v/O=%v/OU=%v/CN=%s", i, subject.Country, subject.Province, subject.Locality, subject.Organization, subject.OrganizationalUnit, subject.CommonName)
		log.Printf("   i:/C=%v/ST=%v/L=%v/O=%v/OU=%v/CN=%s", issuer.Country, issuer.Province, issuer.Locality, issuer.Organization, issuer.OrganizationalUnit, issuer.CommonName)
	}
	log.Print(">>>>>>>>>>>>>>>> State End <<<<<<<<<<<<<<<<")
}

func redirect(w http.ResponseWriter, req *http.Request) {

	host, _, err := net.SplitHostPort(req.Host)
	if err != nil {
		host = req.Host
	}
	host += addr

	targetURL := url.URL{Scheme: "https", Host: host, Path: req.URL.Path, RawQuery: req.URL.RawQuery}
	target := targetURL.String()

	log.Printf("redirect to: %s", target)
	http.Redirect(w, req, target,
		// see @andreiavrammsd comment: often 307 > 301
		http.StatusTemporaryRedirect)
}

// ***************************************************************************

func printConn(conn *tls.Conn) {
	state := conn.ConnectionState()
	printConnState(&state)
}

func main02() {
	flag.Parse()

	// Load CA cert
	caCert, err := ioutil.ReadFile(*caFile)
	if err != nil {
		log.Fatal(err)

	}

	// Load Server cert
	cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
	if err != nil {
		log.Fatal(err)
	}

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)
	if !ok {
		panic("failed to parse CA certificate")
	}

	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,

		Certificates: []tls.Certificate{cert},
	}

	ln, err := tls.Listen("tcp", addr, tlsConfig)
	if err != nil {
		log.Fatal("listen failed: %s", err.Error())
	}
	log.Printf("listen on %s", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("accept failed: %s", err.Error())
			break
		}
		log.Printf("connection open: %s", conn.RemoteAddr())
		printConn(conn.(*tls.Conn))

		go func(c net.Conn) {
			wr, _ := io.Copy(c, c)
			c.Close()
			log.Printf("connection close: %s, written: %d", conn.RemoteAddr(), wr)
		}(conn)
	}

}

func main() {
	main01()
}
