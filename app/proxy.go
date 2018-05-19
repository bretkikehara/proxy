package app

import (
	"github.com/bretkikehara/proxy-test2/proxy"

	"gopkg.in/urfave/cli.v2"
)

var proxyCmd = &cli.Command{
	Name:   "proxy",
	Usage:  "run the proxy server",
	Action: proxyFn,
	Flags:  []cli.Flag{},
}

func proxyFn(ctx *cli.Context) error {
	proxy.New(&proxy.Config{
		Proto: "http",
	})
	return nil
}
