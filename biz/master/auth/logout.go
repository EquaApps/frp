package auth

import (
	"github.com/EquaApps/frp/common"
	"github.com/EquaApps/frp/conf"
	"github.com/EquaApps/frp/pb"
	"github.com/gin-gonic/gin"
)

func RemoveJWTHandler(c *gin.Context) {
	c.SetCookie(conf.Get().App.CookieName,
		"", -1,
		conf.Get().App.CookiePath,
		conf.Get().App.CookieDomain,
		conf.Get().App.CookieSecure,
		conf.Get().App.CookieHTTPOnly)
	common.OKResp(c, &pb.CommonResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS,
		Message: "ok"}})
}
