package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	scmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/peggyjv/sommelier/app"
	"github.com/peggyjv/sommelier/cmd/sommelier/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := scmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)
		default:
			os.Exit(1)
		}
	}
}
