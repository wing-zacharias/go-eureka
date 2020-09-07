package erpc

/**
  @author: wing
  @date: 2020/9/4
  @comment:
**/
import (
	"errors"
	"go-eureka/config"
	"go-eureka/util"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/kataras/golog"
)

/**
* @author: wing
* @time: 2020/9/4 9:38
* @param:
* @return:
* @comment: eureka instance pool
**/
var instancePool []Instance

/**
* @author: wing
* @time: 2020/9/4 9:38
* @param:
* @return:
* @comment: feign client handle
**/
type FeignClient struct {
	App         *Application
	ContextPath string
	logger      *golog.Logger
}

/**
* @author: wing
* @time: 2020/9/4 13:19
* @param:
* @return:
* @comment: create 1 feign client
**/
func NewFeignClient(eurekaService *EurekaService, appName string, contextPath string) *FeignClient {
	app := eurekaService.GetApplication(appName)
	return &FeignClient{
		App:         app,
		ContextPath: contextPath,
		logger:      config.Global.Logger,
	}
}

/**
* @author: wing
* @time: 2020/9/4 13:09
* @param:
* @return:
* @comment: special logger
**/
func (f *FeignClient) SetLogger(logger *golog.Logger) {
	f.logger = logger
}

/**
* @author: wing
* @time: 2020/9/4 9:39
* @param:
* @return:
* @comment: refresh eureka install pool
**/
func (f *FeignClient) refreshInstancePool() {
	for index, lInst := range instancePool {
		instNotInactive := true
		for _, inst := range f.App.Instances {
			if lInst.InstanceId == inst.InstanceId {
				instNotInactive = false
				break
			}
		}
		if instNotInactive {
			instancePool = append(instancePool[:index], instancePool[index+1:]...)
		}
	}
	if len(instancePool) == 0 {
		for _, inst := range f.App.Instances {
			if inst.Status == "UP" {
				instancePool = append(instancePool, inst)
			}
		}
	}
}

/**
* @author: wing
* @time: 2020/9/4 9:39
* @param:
* @return:
* @comment: select 1 instance url for service
**/
func (f *FeignClient) createBaseUrl() string {
	var BaseUrl string
	f.refreshInstancePool()
	if len(instancePool) == 0 {
		f.logger.Errorf("FeignClient.createBaseUrl:cannot find feign instance! ")
		panic(errors.New("cannot find feign instance! "))
	}
	index := rand.Intn(len(instancePool))
	instance := instancePool[index]
	instancePool = append(instancePool[:index], instancePool[index+1:]...)
	ip := instance.IpAddr
	port := instance.Port.Port
	url := "http://" + ip + ":" + strconv.Itoa(port)
	if f.ContextPath != "" {
		fcp := util.FixEndpoint(f.ContextPath)
		BaseUrl = url + fcp
	}
	return BaseUrl
}

/**
* @author: wing
* @time: 2020/9/4 9:40
* @param:
* @return:
* @comment: get this feign client handle
**/
func (f *FeignClient) GetClient() *FeignClient {
	return f
}

/**
* @author: wing
* @time: 2020/9/4 9:40
* @param:
* @return:
* @comment:  set feign context path
**/
func (f *FeignClient) SetContextPath(contextPath string) {
	f.ContextPath = contextPath
}

/**
* @author: wing
* @time: 2020/9/4 9:40
* @param:
* @return:
* @comment: remote call for get
**/
func (f *FeignClient) GetForEntity(endpoint string) ([]byte, error) {
	url := f.createBaseUrl()
	fep := util.FixEndpoint(endpoint)
	url += fep
	f.logger.Debugf("%s:%s", http.MethodGet, url)
	res, err := http.Get(url)
	if err != nil {
		f.logger.Errorf("FeignClient.GetForEntity.1:%s", err)
		return nil, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			f.logger.Errorf("FeignClient.GetForEntity.end:%s", err)
		}
	}()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		f.logger.Errorf("FeignClient.GetForEntity.2:%s", err)
		return nil, err
	}
	return b, nil
}

/**
* @author: wing
* @time: 2020/9/4 9:40
* @param:
* @return:
* @comment: remote call for post
**/
func (f *FeignClient) PostForEntity(endpoint string, contextType string, body io.Reader) ([]byte, error) {
	url := f.createBaseUrl()
	fep := util.FixEndpoint(endpoint)
	url += fep
	f.logger.Debugf("%s:%s", http.MethodPost, url)
	res, err := http.Post(url, contextType, body)
	if err != nil {
		f.logger.Errorf("FeignClient.PostForEntity.1:%s", err)
		return nil, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			f.logger.Errorf("FeignClient.PostForEntity.end:%s", err)
		}
	}()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		f.logger.Errorf("FeignClient.PostForEntity.2:%s", err)
		return nil, err
	}
	return b, nil
}
