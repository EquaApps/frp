package auth

import (
	"github.com/johncoker233/frpaaa/cache"
	"github.com/johncoker233/frpaaa/dao"
	"github.com/johncoker233/frpaaa/models"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

func InitAuth() {
	logrus.Info("start to init frp user auth token")

	u, err := dao.AdminGetAllUsers()
	if err != nil {
		logrus.WithError(err).Fatalf("init frp user auth token failed")
	}

	lo.ForEach(u, func(user *models.UserEntity, _ int) {
		cache.Get().Set([]byte(user.GetUserName()), []byte(user.GetToken()), 0)
	})

	logrus.Infof("init frp user auth token success, count: %d", len(u))
}
