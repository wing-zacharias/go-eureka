package client

/**
  @author: wing
  @date: 2020/9/4
  @comment:
**/

import (
	"go-eureka/erpc"

	"github.com/kataras/iris/v12/mvc"
)

/**
  @author: wing
  @date: 2020/9/5
  @comment:
**/

var mainController *MainController

/**
* @author: wing
* @time: 2020/9/5 14:41
* @param:
* @return:
* @comment: main controller
**/
type MainController struct {
	xxFeignClient *erpc.FeignClient
}

/**
* @author: wing
* @time: 2020/9/5 14:50
* @param:
* @return:
* @comment:
**/
func init() {
	eurekaBase := erpc.GetEurekaBaseFromConfig()
	eurekaService := erpc.NewEurekaService(eurekaBase)
	mainController = &MainController{
		xxFeignClient: eurekaService.GetFeignClient("feign client appName", "/contextPath"),
	}
}

/**
* @author: wing
* @time: 2020/9/5 14:41
* @param:
* @return:
* @comment: router "/${contextPath}/",${contextPath}此处实际为test
**/
func (m *MainController) Get() mvc.Result {
	return mvc.Response{
		ContentType: "text/html",
		Text:        "<h1>Welcome</h1>",
	}
}

/**
* @author: wing
* @time: 2020/9/5 14:44
* @param:
* @return:
* @comment: router "/${contextPath}/rpc",${contextPath}此处实际为test
**/
func (m *MainController) GetRpc() interface{} {
	res, _ := m.xxFeignClient.GetForEntity("/endpoint")
	return res
}
