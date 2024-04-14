package main

import (
	bizserver "github.com/EquaApps/frp/biz/server"
	"github.com/EquaApps/frp/conf"
	"github.com/EquaApps/frp/pb"
	"github.com/EquaApps/frp/rpc"
	"github.com/EquaApps/frp/services/api"
	"github.com/EquaApps/frp/services/rpcclient"
	"github.com/EquaApps/frp/utils"
	"github.com/EquaApps/frp/watcher"
	"github.com/fatedier/golib/crypto"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc"
)

func runServer() {
	var (
		clientID     = conf.Get().Client.ID
		clientSecret = conf.Get().Client.Secret
	)
	crypto.DefaultSalt = conf.Get().App.Secret
	logrus.Infof("start to run server")

	if len(clientID) == 0 {
		logrus.Fatal("client id cannot be empty")
	}

	router := bizserver.NewRouter()
	api.MustInitApiService(conf.ServerAPIListenAddr(), router)

	a := api.GetAPIService()
	defer a.Stop()

	cred, err := utils.TLSClientCertNoValidate(rpc.GetClientCert(clientID, clientSecret, pb.ClientType_CLIENT_TYPE_FRPS))
	if err != nil {
		logrus.Fatal(err)
	}
	conf.ClientCred = cred
	rpcclient.MustInitClientRPCSerivce(
		clientID,
		clientSecret,
		pb.Event_EVENT_REGISTER_SERVER,
		bizserver.HandleServerMessage,
	)

	r := rpcclient.GetClientRPCSerivce()
	defer r.Stop()

	w := watcher.NewClient(bizserver.PullConfig, clientID, clientSecret)
	defer w.Stop()

	initServerOnce(clientID, clientSecret)

	var wg conc.WaitGroup
	wg.Go(r.Run)
	wg.Go(w.Run)
	wg.Go(a.Run)
	wg.Wait()
}

func initServerOnce(clientID, clientSecret string) {
	err := bizserver.PullConfig(clientID, clientSecret)
	if err != nil {
		logrus.WithError(err).Errorf("cannot pull server config, wait for retry")
	}
}
