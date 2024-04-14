package main

import (
	"github.com/johncoker233/frpaaa/conf"
	"github.com/johncoker233/frpaaa/rpc"
	"github.com/spf13/cobra"
)

func main() {
	cobra.MousetrapHelpText = ""

	initLogger()
	initCommand()
	conf.InitConfig()
	rpc.InitRPCClients()

	rootCmd.Execute()
}
