package client

import (
	"context"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/EquaApps/frp/dao"
	"github.com/EquaApps/frp/models"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

func SyncTunnel(ctx context.Context, userInfo models.UserInfo) error {
	clis, err := dao.GetAllClients(userInfo)
	if err != nil {
		return err
	}
	lo.ForEach(clis, func(cli *models.ClientEntity, _ int) {
		cfg, err := cli.GetConfigContent()
		if err != nil {
			logrus.WithError(err).Errorf("cannot get client config content, id: [%s]", cli.ClientID)
			return
		}

		cfg.User = userInfo.GetUserName()
		cfg.Metadatas = map[string]string{
			common.FRPAuthTokenKey: userInfo.GetToken(),
		}
		if err := cli.SetConfigContent(*cfg); err != nil {
			logrus.WithError(err).Errorf("cannot set client config content, id: [%s]", cli.ClientID)
			return
		}

		if err := dao.UpdateClient(userInfo, cli); err != nil {
			logrus.WithError(err).Errorf("cannot update client, id: [%s]", cli.ClientID)
			return
		}
		logrus.Infof("update client success, id: [%s]", cli.ClientID)
	})
	return nil
}
