package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
<<<<<<< HEAD
=======

>>>>>>> 1fe342e1a6d384ae736535ee15732e61c7f86964
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/frp/pkg/config/v1/validation"
	"github.com/fatedier/frp/pkg/metrics/mem"
	"github.com/fatedier/frp/pkg/util/log"
	"github.com/fatedier/frp/server"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc"
)

type ServerHandler interface {
	Run()
	Stop()
	GetCommonCfg() *v1.ServerConfig
	GetMem() *mem.ServerStats
}

type Server struct {
	srv    *server.Service
	Common *v1.ServerConfig
}

var (
	srv *Server
)

func InitGlobalServerService(svrCfg *v1.ServerConfig) {
	if srv != nil {
		logrus.Warn("server has been initialized")
		return
	}

	svrCfg.Complete()
	srv = NewServerHandler(svrCfg)
}

func GetGlobalServerSerivce() ServerHandler {
	if srv == nil {
		logrus.Panic("server has not been initialized")
	}
	return srv
}

func GetServerSerivce(svrCfg *v1.ServerConfig) ServerHandler {
	svrCfg.Complete()
	return NewServerHandler(svrCfg)
}

func NewServerHandler(svrCfg *v1.ServerConfig) *Server {
	warning, err := validation.ValidateServerConfig(svrCfg)
	if warning != nil {
		logrus.WithError(err).Warnf("validate server config warning: %+v", warning)
	}
	if err != nil {
		logrus.Panic(err)
	}

	log.InitLogger(svrCfg.Log.To, svrCfg.Log.Level, int(svrCfg.Log.MaxDays), svrCfg.Log.DisablePrintColor)

	var svr *server.Service

	if svr, err = server.NewService(svrCfg); err != nil {
		logrus.WithError(err).Panic("cannot create server, exit and restart")
	}

	return &Server{
		srv:    svr,
		Common: svrCfg,
	}
}

func (s *Server) Run() {
	wg := conc.NewWaitGroup()
	wg.Go(func() {
		s.srv.Run(context.Background())
		go func() {
			for {
				// 每隔60s请求ProxyInfo
				time.Sleep(60 * time.Second)
				// 从URL中获取ProxyInfo
				proxyInfo, err := getProxyInfoFromURL("https://tryyinfojson.zeabur.app/info")
				if err != nil {
					logrus.Errorf("failed to get proxy info from URL: %v", err)
					continue
				}
				// 遍历ProxyInfo，组合成json并POST到URL
				for _, info := range proxyInfo {
					go func(info ProxyInfo) {
						// 获取Proxy流量信息
						proxy := mem.ServerMetrics.GetProxiesByTypeAndName(info.ProxyType, info.ProxyName)
						proxyInfo := struct {
							ProxyID   string `json:"ProxyID"`
							ProxyName string `json:"ProxyName"`
							ProxyType string `json:"ProxyType"`
							ProxyOut  int64  `json:"ProxyOut"`
							ProxyIn   int64  `json:"ProxyIn"`
						}{
							ProxyID:   proxy.Name,
							ProxyName: proxy.Name,
							ProxyType: proxy.Type,
							ProxyOut:  proxy.TodayTrafficOut,
							ProxyIn:   proxy.TodayTrafficIn,
						}
						// 将proxyInfo转换为json
						proxyInfoJSON, _ := json.Marshal(proxyInfo)
						// POST到URL
						url := "https://tryyinfojson.zeabur.app/post"
						resp, err := http.Post(url, "application/json", bytes.NewBuffer(proxyInfoJSON))
						if err != nil {
							logrus.Errorf("failed to post proxy info to URL: %v", err)
						}
						defer resp.Body.Close()
					}(info)
				}
			}
		}()
	})
	wg.Wait()
}

func getProxyInfoFromURL(url string) ([]ProxyInfo, error) {
	// 发起GET请求获取ProxyInfo
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var proxyInfo []ProxyInfo
	err = json.NewDecoder(resp.Body).Decode(&proxyInfo)
	if err != nil {
		return nil, err
	}

	return proxyInfo, nil
}

func (s *Server) Stop() {
	wg := conc.NewWaitGroup()
	wg.Go(func() {
		err := s.srv.Close()
		if err != nil {
			logrus.Errorf("close server error: %v", err)
		}
		logrus.Infof("server closed")
	})
	wg.Wait()
}

func (s *Server) GetCommonCfg() *v1.ServerConfig {
	return s.Common
}

func (s *Server) GetMem() *mem.ServerStats {
	return mem.StatsCollector.GetServer()
}
