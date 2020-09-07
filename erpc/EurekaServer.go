package erpc

/**
  @author: wing
  @date: 2020/9/4
  @comment:
**/
import (
	"context"
	"go-eureka/config"
	"go-eureka/network"
	"go-eureka/util"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/kataras/golog"
)

/**
* @author: wing
* @time: 2020/9/4 9:27
* @param:
* @return:
* @comment: EurekaServer entity
**/
type EurekaServer struct {
	Cluster      *Cluster
	HttpClient   *http.Client
	HttpResponse *http.Response
	logger       *golog.Logger
}

/**
* @author: wing
* @time: 2020/9/4 9:27
* @param: eurekaNodes:at lease 1
* @return: EurekaServer
* @comment: create EurekaServer from eurekaNodes
**/
func GetEurekaBase(eurekaNodes []string) *EurekaServer {
	logger := config.Global.Logger
	if len(eurekaNodes) != 0 {
		return &EurekaServer{
			Cluster:    GetCluster(eurekaNodes),
			HttpClient: &http.Client{},
			logger:     logger,
		}
	} else {
		logger.Errorf("Can not find any eureka server! ")
	}
	return nil
}

/**
* @author: wing
* @time: 2020/9/4 11:40
* @param:
* @return:
* @comment: read config from file
**/
func GetEurekaBaseFromConfig() *EurekaServer {
	conf := config.Global.EurekaConfig
	logger := config.Global.Logger
	if conf.AutoDetect.Allow {
		return EurekaBaseAutoDetect()
	}
	if len(conf.EurekaNodes) != 0 {
		return GetEurekaBase(conf.EurekaNodes)
	} else {
		logger.Errorf("Can not find any eureka server! ")
	}
	return nil
}

/**
* @author: wing
* @time: 2020/9/4 10:45
* @param:
* @return:
* @comment: auto detect local network eureka server
**/
func EurekaBaseAutoDetect() *EurekaServer {
	conf := config.Global.EurekaConfig
	logger := config.Global.Logger
	var eurekaNodes []string
	ip, maskLen, err := util.CidrDisassemble(conf.AutoDetect.Network)
	if err != nil {
		ip, maskLen = util.GetLocalIp()
	}
	servers := network.PortOpenFullScan(ip, maskLen, conf.AutoDetect.Port, conf.AutoDetect.Timeout)
	for _, server := range servers {
		eurekaNode := "http://" + server.Ip + ":" + server.Port + util.FixEndpoint(conf.AutoDetect.ContextPath)
		logger.Debugf("Detect eureka server: %s", eurekaNode)
		eurekaNodes = append(eurekaNodes, eurekaNode)
	}
	if len(eurekaNodes) == 0 {
		logger.Errorf("Can not detect any eureka server! ")
		return nil
	}
	return GetEurekaBase(eurekaNodes)
}

/**
* @author: wing
* @time: 2020/9/4 13:06
* @param:
* @return:
* @comment: special logger
**/
func (e *EurekaServer) SetLogger(logger *golog.Logger) {
	e.logger = logger
}

/**
* @author: wing
* @time: 2020/9/4 9:28
* @param:
* @return: eureka url
* @comment: select cluster leader for service
**/
func (e *EurekaServer) GetEurekaBaseUrl() string {
	eurekaBaseUrl := e.Cluster.Leader
	return eurekaBaseUrl
}

/**
* @author: wing
* @time: 2020/9/4 9:29
* @param: client:set http client for service if need
* @return:
* @comment: create http client for service
**/
func (e *EurekaServer) SetHttpClient(client *http.Client) {
	e.HttpClient = client
}

/**
* @author: wing
* @time: 2020/9/4 9:30
* @param: endpoint: get info from eureka
* @return:
* @comment: get service
**/
func (e *EurekaServer) Get(endpoint string) ([]byte, error) {
	reqUrl := e.GetEurekaBaseUrl() + util.FixEndpoint(endpoint)
	e.logger.Debugf("EurekaServer.Get:%s-%s", http.MethodGet, reqUrl)
	res, err := http.Get(reqUrl)
	if err != nil {
		e.logger.Errorf("EurekaServer.Get.1:%s", err)
		return nil, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			e.logger.Errorf("EurekaServer.Get.end:%s", err)
		}
	}()
	e.HttpResponse = res
	if res.StatusCode == http.StatusTemporaryRedirect {
		resUrl, err := res.Location()
		if err == nil {
			e.Cluster.updateLeaderByRequestUrl(resUrl)
		} else {
			e.logger.Errorf("EurekaServer.Get.2:%s", err)
		}
	}
	resByte, err := ioutil.ReadAll(res.Body)
	if err != nil {
		e.logger.Errorf("EurekaServer.Get.3:%s", err)
		return nil, err
	}
	return resByte, nil
}

/**
* @author: wing
* @time: 2020/9/4 9:32
* @param: post info to eureka
* @return:
* @comment: post service
**/
func (e *EurekaServer) Post(endpoint string, contentType string, body io.Reader) ([]byte, error) {
	reqUrl := e.GetEurekaBaseUrl() + util.FixEndpoint(endpoint)
	e.logger.Debugf("EurekaServer.Post:%s-%s", http.MethodPost, reqUrl)
	res, err := http.Post(reqUrl, contentType, body)
	if err != nil {
		e.logger.Errorf("EurekaServer.Post.1:%s", err)
		return nil, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			e.logger.Errorf("EurekaServer.Post.end:%s", err)
		}
	}()
	e.HttpResponse = res
	if res.StatusCode == http.StatusTemporaryRedirect {
		resUrl, err := res.Location()
		if err == nil {
			e.Cluster.updateLeaderByRequestUrl(resUrl)
		} else {
			e.logger.Errorf("EurekaServer.Post.2:%s", err)
		}
	}
	resByte, err := ioutil.ReadAll(res.Body)
	if err != nil {
		e.logger.Errorf("EurekaServer.Post.3:%s", err)
		return nil, err
	}
	return resByte, nil
}

/**
* @author: wing
* @time: 2020/9/4 9:33
* @param: put info to eureka
* @return:
* @comment: put service
**/
func (e *EurekaServer) Put(endpoint string, body io.Reader) ([]byte, error) {
	reqUrl := e.GetEurekaBaseUrl() + util.FixEndpoint(endpoint)
	e.logger.Debugf("EurekaServer.Put:%s-%s", http.MethodPut, reqUrl)
	req, err := http.NewRequest(http.MethodPut, reqUrl, body)
	if err != nil {
		e.logger.Errorf("EurekaServer.Put.1:%s", err)
		return nil, err
	}
	res, err := e.HttpClient.Do(req)
	if err != nil {
		e.logger.Errorf("EurekaServer.Put.2:%s", err)
		return nil, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			e.logger.Errorf("EurekaServer.Put.end:%s", err)
		}
	}()
	e.HttpResponse = res
	if res.StatusCode == http.StatusTemporaryRedirect {
		resUrl, err := res.Location()
		if err == nil {
			e.Cluster.updateLeaderByRequestUrl(resUrl)
		} else {
			e.logger.Errorf("EurekaServer.Put.3:%s", err)
		}
	}
	resByte, err := ioutil.ReadAll(res.Body)
	if err != nil {
		e.logger.Errorf("EurekaServer.Put.4:%s", err)
		return nil, err
	}
	return resByte, nil
}

/**
* @author: wing
* @time: 2020/9/4 9:33
* @param: delete info to eureka
* @return:
* @comment: delete service
**/
func (e *EurekaServer) Delete(endpoint string, body io.Reader) ([]byte, error) {
	reqUrl := e.GetEurekaBaseUrl() + util.FixEndpoint(endpoint)
	e.logger.Debugf("EurekaServer.Delete:%s-%s", http.MethodDelete, reqUrl)
	req, err := http.NewRequest(http.MethodDelete, reqUrl, body)
	if err != nil {
		e.logger.Errorf("EurekaServer.Delete.1:%s", err)
		return nil, err
	}
	res, err := e.HttpClient.Do(req)
	if err != nil {
		e.logger.Errorf("EurekaServer.Delete.2:%s", err)
		return nil, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			e.logger.Errorf("EurekaServer.Delete.end:%s", err)
		}
	}()
	e.HttpResponse = res
	if res.StatusCode == http.StatusTemporaryRedirect {
		resUrl, err := res.Location()
		if err == nil {
			e.Cluster.updateLeaderByRequestUrl(resUrl)
		} else {
			e.logger.Errorf("EurekaServer.Delete.3:%s", err)
		}
	}
	resByte, err := ioutil.ReadAll(res.Body)
	if err != nil {
		e.logger.Errorf("EurekaServer.Delete.4:%s", err)
		return nil, err
	}
	return resByte, nil
}

/**
* @author: wing
* @time: 2020/9/4 9:34
* @param:
* @return:
* @comment: make special http request,may be cancelable
**/
func (e *EurekaServer) DoRequest(request *http.Request, resChan chan<- []byte) {
	e.logger.Debugf("EurekaServer.DoRequest:%s-%s", request.Method, request.URL.Path)
	go func() {
		res, err := e.HttpClient.Do(request)
		if err != nil {
			e.logger.Errorf("EurekaServer.DoRequest.1:%s", err)
			resChan <- []byte(err.Error())
		}
		defer func() {
			if err := res.Body.Close(); err != nil {
				e.logger.Errorf("EurekaServer.DoRequest.end:%s", err)
			}
		}()
		e.HttpResponse = res
		if res.StatusCode == http.StatusTemporaryRedirect {
			resUrl, err := res.Location()
			if err == nil {
				e.Cluster.updateLeaderByRequestUrl(resUrl)
			} else {
				e.logger.Errorf("EurekaServer.DoRequest.2:%s", err)
			}
		}
		resByte, err := ioutil.ReadAll(res.Body)
		if err != nil {
			e.logger.Errorf("EurekaServer.DoRequest.3:%s", err)
			resChan <- []byte(err.Error())
		}
		resChan <- resByte
	}()
}

/**
* @author: wing
* @time: 2020/9/4 9:35
* @param:
* @return:
* @comment: cancel http request
**/
func (e *EurekaServer) CancelRequest(request *http.Request, resChan chan<- []byte) {
	_, cancel := context.WithCancel(request.Context())
	defer cancel()
	resChan <- []byte("request canceled! ")
}
