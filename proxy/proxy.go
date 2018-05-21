package proxy

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func handleTunneling(w http.ResponseWriter, r *http.Request) {
	dest_conn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	client_conn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	go transfer(dest_conn, client_conn)
	go transfer(client_conn, dest_conn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func handleHTTPReq(w http.ResponseWriter, req *http.Request) {
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	var bd []byte
	bd, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	if resp.Header.Get("Content-Encoding") == "gzip" {
		bd, err = UnGzip(bd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
	}
	bd = bytes.Join([][]byte{[]byte("hello"), bd}, []byte("\n"))

	if resp.Header.Get("Content-Encoding") == "gzip" {
		bd, err = Gzip(bd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
	}

	resp.Header.Set("Content-Length", fmt.Sprintf("%d", len(bd)))
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	w.Write(bd)
}

func mutate(w http.ResponseWriter, req *http.Request, resp *http.Response, body []byte) (mBody []byte, mHeader http.Header, status int) {
	return body, http.Header{}, resp.StatusCode
}

func handleRequest(req *http.Request) (resp *http.Response, bd []byte, mStatus int, err error) {
	resp, err = http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return resp, nil, http.StatusServiceUnavailable, err
	}
	defer resp.Body.Close()

	bd, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, resp.StatusCode, err
	}
	return resp, bd, resp.StatusCode, nil
}

func handleHTTP(cfg *Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodConnect {
			fmt.Printf("[CON] %s\n", req.URL.String())
			handleTunneling(w, req)
		} else {
			fmt.Printf("[REQ] %s\n", req.URL.String())
			handleHTTPReq(w, req)
			// var resp *http.Response
			// var mBody []byte
			// var mHeader http.Header
			// var mStatus int
			// var err error
			// resp, mBody, mStatus, err = handleRequest(req)
			// if err == nil {
			// 	mBody, mHeader, mStatus = mutate(w, req, resp, mBody)
			// } else {
			// 	mHeader = resp.Header
			// 	mStatus = resp.StatusCode
			// }
			// copyHeader(w.Header(), mHeader)
			// w.WriteHeader(mStatus)
			// fmt.Fprint(w, mBody)
		}
	}
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
	Port   string
	SSLPem string
	SSLKey string
}

// New creates the new proxy
func New(cfg *Config) func() {
	fmt.Printf("starting proxy server on :%s\n", cfg.Port)
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: http.HandlerFunc(handleHTTP(cfg)),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	go func() {
		var err error
		if cfg.Proto == "https" {
			err = server.ListenAndServeTLS(cfg.SSLPem, cfg.SSLKey)
		} else {
			err = server.ListenAndServe()
		}
		if err != nil {
			fmt.Printf("proxy server failed: %s\n", err)
		}
	}()
	return func() {
		server.Close()
	}
}
