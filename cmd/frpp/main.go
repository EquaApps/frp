package main

import (
	"github.com/johncoker233/frpaaa/conf"
	"github.com/johncoker233/frpaaa/rpc"
)

func main() {
	initLogger()
	initCommand()
	conf.InitConfig()
	rpc.InitRPCClients()

	rootCmd.Execute()
}
