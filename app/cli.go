package app

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/bretkikehara/proxy/proxy"

	"gopkg.in/urfave/cli.v2"
)

var cliCmd = &cli.Command{
	Name:   "cli",
	Usage:  "run the proxy server",
	Action: cliFn,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Value:   "8889",
			Usage:   "proxy server port",
		},
		&cli.StringFlag{
			Name:    "url",
			Aliases: []string{"u"},
			Value:   "http://example.com",
			Usage:   "URL to retrieve",
		},
	},
}

func cliFn(ctx *cli.Context) error {
	var wg sync.WaitGroup
	var closePxy func()
	var err error

	wg.Add(1)
	go func() {
		closePxy = proxy.New(&proxy.Config{
			Proto: "http",
		})
		fmt.Printf("proxy server ready\n")
		defer wg.Done()
	}()
	wg.Wait()
	if err != nil {
		return err
	}
	if closePxy != nil {
		defer closePxy()
	}
	return fetchURL(ctx.String("url"), "127.0.0.1:"+ctx.String("port"))
}

func fetchURL(ur string, pxyServer string) error {
	fmt.Printf("Making request\n")
	c := &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			Proxy: http.ProxyURL(&url.URL{
				Host: pxyServer,
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

	fmt.Printf("Parsing the resp\n")
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("%s", body)
	return nil
}
