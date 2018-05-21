package app

import (
	"os"
	"os/signal"
	"sync"

	"github.com/bretkikehara/proxy/proxy"

	"gopkg.in/urfave/cli.v2"
)

var serveCmd = &cli.Command{
	Name:   "serve",
	Usage:  "run the proxy server",
	Action: serveFn,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Value:   "8889",
			Usage:   "proxy server port",
		},
	},
}

func WaitForCtrlC() {
	var end_waiter sync.WaitGroup
	end_waiter.Add(1)
	var signal_channel chan os.Signal
	signal_channel = make(chan os.Signal, 1)
	signal.Notify(signal_channel, os.Interrupt)
	go func() {
		<-signal_channel
		end_waiter.Done()
	}()
	end_waiter.Wait()
}

func serveFn(ctx *cli.Context) error {
	close := proxy.New(&proxy.Config{
		Port: ctx.String("port"),
	})
	defer close()
	WaitForCtrlC()
	return nil
}
