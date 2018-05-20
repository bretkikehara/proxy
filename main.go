package main

import (
	"fmt"
	"os"

	"github.com/bretkikehara/proxy/app"
)

func main() {
	if err := app.App.Run(os.Args); err != nil {
		fmt.Printf("error running the app: %s\n", err)
	}
}
