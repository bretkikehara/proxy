package app

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/bretkikehara/proxy-test2/proxy"

	"gopkg.in/urfave/cli.v2"
)

var testCmd = &cli.Command{
	Name:   "test",
	Usage:  "run the proxy server",
	Action: testFn,
	Flags:  []cli.Flag{},
}

func testFn(ctx *cli.Context) error {
	var wg sync.WaitGroup
	var closePxy func()
	var err error

	wg.Add(1)
	go func() {
		closePxy, err = proxy.New(&proxy.Config{
			Proto: "http",
		})
		defer wg.Done()
	}()
	wg.Wait()
	if err != nil {
		return err
	}
	if closePxy != nil {
		defer closePxy()
	}
	return fetchURL("http://example.com")
}

func fetchURL(ur string) error {
	fmt.Printf("Making request\n")
	c := &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			Proxy: http.ProxyURL(&url.URL{
				Host: "127.0.0.1:8888",
			}),
		},
	}

	var resp *http.Response
	var err error
	resp, err = c.Get(ur)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("%s", body)
	return nil
}
