package main

import (
	bizclient "github.com/johncoker233/frpaaa/biz/client"
	"github.com/johncoker233/frpaaa/conf"
	"github.com/johncoker233/frpaaa/pb"
	"github.com/johncoker233/frpaaa/rpc"
	"github.com/johncoker233/frpaaa/services/rpcclient"
	"github.com/johncoker233/frpaaa/utils"
	"github.com/johncoker233/frpaaa/watcher"
	"github.com/fatedier/golib/crypto"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc"
)

func runClient() {
	var (
		clientID     = conf.Get().Client.ID
		clientSecret = conf.Get().Client.Secret
	)
	crypto.DefaultSalt = conf.Get().App.Secret
	logrus.Infof("start to run client")
	if len(clientSecret) == 0 {
		logrus.Fatal("client secret cannot be empty")
	}

	if len(clientID) == 0 {
		logrus.Fatal("client id cannot be empty")
	}

	cred, err := utils.TLSClientCertNoValidate(rpc.GetClientCert(clientID, clientSecret, pb.ClientType_CLIENT_TYPE_FRPC))
	if err != nil {
		logrus.Fatal(err)
	}
	conf.ClientCred = cred

	rpcclient.MustInitClientRPCSerivce(
		clientID,
		clientSecret,
		pb.Event_EVENT_REGISTER_CLIENT,
		bizclient.HandleServerMessage,
	)
	r := rpcclient.GetClientRPCSerivce()
	defer r.Stop()

	w := watcher.NewClient(bizclient.PullConfig, clientID, clientSecret)
	defer w.Stop()

	initClientOnce(clientID, clientSecret)

	var wg conc.WaitGroup
	wg.Go(r.Run)
	wg.Go(w.Run)
	wg.Wait()
}

func initClientOnce(clientID, clientSecret string) {
	err := bizclient.PullConfig(clientID, clientSecret)
	if err != nil {
		logrus.WithError(err).Errorf("cannot pull client config, wait for retry")
	}
}
