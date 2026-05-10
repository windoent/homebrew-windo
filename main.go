package main

import (
	"os"

	"github.com/windoent/homebrew-windo/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
