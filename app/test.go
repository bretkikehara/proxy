package app

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"gopkg.in/urfave/cli.v2"
)

var testCmd = &cli.Command{
	Name:   "test",
	Usage:  "run the proxy server",
	Action: testFn,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Value:   "8889",
			Usage:   "proxy server port",
		},
	},
}

func testFn(ctx *cli.Context) error {
	u, err := url.Parse("http://localhost:8889")
	if err != nil {
		return err
	}
	tr := &http.Transport{
		Proxy: http.ProxyURL(u),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get("https://www.example.com")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var bd []byte
	bd, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("%s", bd)

	return nil
}
