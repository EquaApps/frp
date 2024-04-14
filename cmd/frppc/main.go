package main

import (
	"github.com/EquaApps/frp/conf"
	"github.com/EquaApps/frp/rpc"
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
