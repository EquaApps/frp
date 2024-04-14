package main

import (
	"github.com/EquaApps/frp/conf"
	"github.com/EquaApps/frp/rpc"
)

func main() {
	initLogger()
	initCommand()
	conf.InitConfig()
	rpc.InitRPCClients()

	rootCmd.Execute()
}
