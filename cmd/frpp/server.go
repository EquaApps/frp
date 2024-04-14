package main

import (
	bizserver "github.com/johncoker233/frpaaa/biz/server"
	"github.com/johncoker233/frpaaa/conf"
	"github.com/johncoker233/frpaaa/pb"
	"github.com/johncoker233/frpaaa/rpc"
	"github.com/johncoker233/frpaaa/services/api"
	"github.com/johncoker233/frpaaa/services/rpcclient"
	"github.com/johncoker233/frpaaa/utils"
	"github.com/johncoker233/frpaaa/watcher"
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
