package client

/**
  @author: wing
  @date: 2020/9/4
  @comment:
**/

import (
	"go-eureka/config"
	"go-eureka/erpc"
	"go-eureka/util"
	"strconv"
	"sync"

	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/mvc"
)

/**
  @author: wing
  @date: 2020/9/5
  @comment:
**/

/**
* @author: wing
* @time: 2020/9/5 14:39
* @param:
* @return:
* @comment: webserver单例
**/
var (
	once sync.Once
	web  *WebServer
)

/**
* @author: wing
* @time: 2020/9/5 14:39
* @param:
* @return:
* @comment: webserver entity
**/
type WebServer struct {
	App      *iris.Application
	AppName  string
	Port     string
	es       *erpc.EurekaService
	ch       chan bool
	wg       *sync.WaitGroup
	logger   *golog.Logger
	instance *erpc.Instance
}

/**
* @author: wing
* @time: 2020/9/5 14:40
* @param:
* @return:
* @comment: create 1 webserver instance
**/
func GetServer(signalChan chan bool) *WebServer {
	once.Do(func() {
		erkSvr := erpc.GetEurekaBaseFromConfig()
		erkSvc := erpc.NewEurekaService(erkSvr)
		appName := config.Global.GlobalConfig.Common.AppName
		port := config.Global.GlobalConfig.Common.Port
		localIp, _ := util.GetLocalIp()
		iPort, _ := strconv.Atoi(port)
		instanceId := localIp + ":" + appName + ":" + config.Global.GlobalConfig.Common.Port
		leaseInfo := &erpc.LeaseInfo{
			EvictionDurationInSecs: uint(30),
		}
		dataCenterInfo := &erpc.DataCenterInfo{
			Name:     "MyOwn",
			Class:    "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo",
			Metadata: nil,
		}
		ePort := &erpc.Port{
			Port:    iPort,
			Enabled: true,
		}
		inst := &erpc.Instance{
			HostName:         util.GetHostname(),
			App:              appName,
			InstanceId:       instanceId,
			IpAddr:           localIp,
			Port:             ePort,
			LeaseInfo:        leaseInfo,
			Status:           "UP",
			DataCenterInfo:   dataCenterInfo,
			VipAddress:       appName,
			SecureVipAddress: "http://" + localIp + ":" + port,
			StatusPageUrl:    "http://" + localIp + ":" + port + "/info",
			HealthCheckUrl:   "http://" + localIp + ":" + port + "/health",
			HomePageUrl:      "http://" + localIp + ":" + port + "/",
		}
		web = &WebServer{
			App:      iris.New(),
			AppName:  appName,
			Port:     port,
			es:       erkSvc,
			ch:       signalChan,
			wg:       &sync.WaitGroup{},
			logger:   config.Global.Logger,
			instance: inst,
		}
		web.App.Use(recover.New())
		web.App.Use(logger.New())
		web.App.Logger().SetLevel(config.Global.GlobalConfig.Log.Console.LogLevel)
		web.App.Get("/", func(ctx iris.Context) {
			ctx.Redirect("/test")
		})
		mvc.New(web.App.Party(config.Global.GlobalConfig.Common.ContextPath)).Handle(mainController)
	})
	return web
}

/**
* @author: wing
* @time: 2020/9/5 14:40
* @param:
* @return:
* @comment: serving
**/
func (s *WebServer) Serving() {
	web.wg.Add(1)
	defer web.wg.Done()
	addr := config.Global.GlobalConfig.Common.Listen + ":" + config.Global.GlobalConfig.Common.Port
	s.logger.Infof("webserver running ...")
	if config.Global.GlobalConfig.Common.AutoRegister {
		if err := s.es.Register(s.AppName, s.instance); err != nil {
			s.logger.Errorf("%s", err)
		}
		s.es.SendHeartbeat(s.AppName, s.instance, 30, s.ch)
		iris.RegisterOnInterrupt(func() {
			s.logger.Warnf("receive exit signal ...")
			s.logger.Warnf("transport exit signal to heartbeat ...")
			s.ch <- true
			s.logger.Warnf("unregister service ...")
			s.es.UnRegister(s.AppName, s.instance.InstanceId)
		})
	}
	web.App.Run(iris.Addr(addr), iris.WithoutServerError(iris.ErrServerClosed), iris.WithCharset("utf-8"))
}
