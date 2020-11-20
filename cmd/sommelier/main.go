package main

import (
	"os"

	"github.com/peggyjv/sommelier/cmd/sommelier/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := cmd.Execute(rootCmd); err != nil {
		os.Exit(1)
	}
}
