package main

import (
	"os"

	"github.com/Alvkoen/barely-incharge/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
