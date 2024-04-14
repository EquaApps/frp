package auth

import (
	"github.com/johncoker233/frpaaa/common"
	"github.com/johncoker233/frpaaa/conf"
	"github.com/johncoker233/frpaaa/pb"
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
