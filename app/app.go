package app

import "gopkg.in/urfave/cli.v2"

// App defines the application
var App = &cli.App{
	Name:  "proxy-test",
	Usage: "Proxies http requests",
	Authors: []*cli.Author{
		{
			Name:  "Bret K. Ikehara",
			Email: "bret.k.ikehara@gmail.com",
		},
	},
	Version: "1.0.0",
	Commands: []*cli.Command{
		cliCmd,
		serveCmd,
		testCmd,
	},
}
