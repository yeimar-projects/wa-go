package main

import (
	"fmt"
	"os"

	"github.com/yeimar-sandbox/wa-go/bootstrap"
)

func main() {
	app := bootstrap.Boot()

	if err := bootstrap.ValidateEnv(); err != nil {
		fmt.Fprintf(os.Stderr, "startup aborted: %v\n", err)
		os.Exit(1)
	}

	app.Start()
}
