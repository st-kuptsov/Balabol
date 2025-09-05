package main

import (
	"fmt"
	"os"

	"github.com/st-kuptsov/balabol/internal/app"
)

var Version = "dev"

func main() {
	if err := app.Run(Version); err != nil {
		fmt.Fprintf(os.Stderr, "application failed: %v\n", err)
		os.Exit(1)
	}
}
