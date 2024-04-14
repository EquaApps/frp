package middleware

import (
	"errors"
	"strings"
	"time"

	"github.com/johncoker233/frpaaa/common"
	"github.com/johncoker233/frpaaa/conf"
	"github.com/johncoker233/frpaaa/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

// JWTMAuth check if authed and set uid to context
func JWTAuth(c *gin.Context) {
	defer func() {
		logrus.WithContext(c).Info("finish jwt middleware")
	}()

	var tokenStr string

	if tokenStr = c.Copy().Query(common.TokenKey); len(tokenStr) != 0 {
		if t, err := utils.ParseToken(conf.JWTSecret(), tokenStr); err == nil {
			for k, v := range t {
				c.Set(k, v)
			}
			logrus.WithContext(c).Infof("query auth success")
			if err = resignAndPatchCtxJWT(c, t, tokenStr); err != nil {
				logrus.WithContext(c).WithError(err).Errorf("resign jwt error")
				common.ErrUnAuthorized(c, "resign jwt error")
				c.Abort()
				return
			}
			c.Next()
			SetToken(c, utils.ToStr(t[common.UserIDKey]))
			return
		}
		logrus.WithContext(c).Infof("query auth failed")
	}

	cookieToken, err := c.Cookie(conf.Get().App.CookieName)
	if err == nil {
		if t, err := utils.ParseToken(conf.JWTSecret(), cookieToken); err == nil {
			for k, v := range t {
				c.Set(k, v)
			}
			logrus.WithContext(c).Infof("cookie auth success")
			if err = resignAndPatchCtxJWT(c, t, cookieToken); err != nil {
				logrus.WithContext(c).WithError(err).Errorf("resign jwt error")
				common.ErrUnAuthorized(c, "resign jwt error")
				c.Abort()
				return
			}
			c.Next()
			return
		} else {
			logrus.WithContext(c).WithError(err).Errorf("jwt middleware parse token error")
			common.ErrUnAuthorized(c, "invalid authorization")
			c.Abort()
			return
		}
	}

	tokenStr = c.Request.Header.Get(common.AuthorizationKey)
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	if tokenStr == "" || tokenStr == "null" {
		logrus.WithContext(c).WithError(errors.New("authorization is empty")).Infof("authorization is empty")
		common.ErrUnAuthorized(c, "invalid authorization")
		c.Abort()
		return
	}

	if t, err := utils.ParseToken(conf.JWTSecret(), tokenStr); err == nil {
		for k, v := range t {
			c.Set(k, v)
		}
		logrus.WithContext(c).Infof("header auth success")
		if err = resignAndPatchCtxJWT(c, t, tokenStr); err != nil {
			logrus.WithContext(c).WithError(err).Errorf("resign jwt error")
			common.ErrUnAuthorized(c, "resign jwt error")
			c.Abort()
			return
		}
		c.Next()
		return
	} else {
		logrus.WithContext(c).WithError(err).Errorf("jwt middleware parse token error")
	}
}

func resignAndPatchCtxJWT(c *gin.Context, t jwt.MapClaims, tokenStr string) error {
	tokenExpire, _ := t.GetExpirationTime()
	now := time.Now().Add(time.Duration(conf.Get().App.CookieAge/2) * time.Second)
	if now.Before(tokenExpire.Time) {
		logrus.WithContext(c).Infof("jwt not going to expire, continue to use old one")
		c.Set(common.TokenKey, tokenStr)
		return nil
	}

	token, err := utils.GetJwtTokenFromMap(conf.JWTSecret(),
		time.Now().Unix(),
		int64(conf.Get().App.CookieAge),
		map[string]string{common.UserIDKey: utils.ToStr(t[common.UserIDKey])})
	if err != nil {
		c.Set(common.TokenKey, tokenStr)
		logrus.WithContext(c).WithError(err).Errorf("resign jwt error")
		return err
	}

	logrus.WithContext(c).Infof("jwt going to expire, resigning token")
	c.Header(common.SetAuthorizationKey, token)
	c.SetCookie(conf.Get().App.CookieName,
		token,
		conf.Get().App.CookieAge,
		conf.Get().App.CookiePath,
		conf.Get().App.CookieDomain,
		conf.Get().App.CookieSecure,
		conf.Get().App.CookieHTTPOnly)
	c.Set(common.TokenKey, token)
	return nil
}

func SetToken(c *gin.Context, uid string) (string, error) {
	logrus.WithContext(c).Infof("set token for uid:[%s]", uid)
	token, err := utils.GetJwtTokenFromMap(conf.JWTSecret(),
		time.Now().Unix(),
		int64(conf.Get().App.CookieAge),
		map[string]string{common.UserIDKey: uid})
	if err != nil {
		return "", err
	}
	c.SetCookie(conf.Get().App.CookieName,
		token,
		conf.Get().App.CookieAge,
		conf.Get().App.CookiePath,
		conf.Get().App.CookieDomain,
		conf.Get().App.CookieSecure,
		conf.Get().App.CookieHTTPOnly)
	c.Set(common.TokenKey, token)
	c.Header(common.SetAuthorizationKey, token)
	return token, nil
}

// PushTokenStr 推送token到客户端
func PushTokenStr(c *gin.Context, tokenStr string) {
	logrus.WithContext(c).Infof("push new token to client")
	c.SetCookie(conf.Get().App.CookieName,
		tokenStr,
		conf.Get().App.CookieAge,
		conf.Get().App.CookiePath,
		conf.Get().App.CookieDomain,
		conf.Get().App.CookieSecure,
		conf.Get().App.CookieHTTPOnly)
	c.Header(common.SetAuthorizationKey, tokenStr)
}
