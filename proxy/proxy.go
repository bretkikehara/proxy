package proxy

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
)

func handleHTTP(w http.ResponseWriter, req *http.Request) {
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

// Config defines the proxy config
type Config struct {
	Proto  string
	SSLPem string
	SSLKey string
}

// New creates the new proxy
func New(cfg *Config) func() {
	proto := cfg.Proto
	if proto != "http" && proto != "https" {
		fmt.Printf("[PXY] proto '%s' not recognized, using http", cfg.Proto)
		proto = "http"
	}
	fmt.Printf("starting proxy server\n")
	server := &http.Server{
		Addr:    ":8889",
		Handler: http.HandlerFunc(handleHTTP),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	go func() {
		var err error
		if proto == "http" {
			err = server.ListenAndServe()
		} else {
			err = server.ListenAndServeTLS(cfg.SSLPem, cfg.SSLKey)
		}
		if err != nil {
			fmt.Printf("proxy server failed:\n%s", err)
		}
	}()
	return func() {
		server.Close()
	}
}
