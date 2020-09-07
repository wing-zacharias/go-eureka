package erpc

/**
  @author: wing
  @date: 2020/9/4
  @comment:
**/
import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"go-eureka/config"
	"net/http"
	"time"

	"github.com/kataras/golog"
)

/**
* @author: wing
* @time: 2020/9/4 13:10
* @param:
* @return:
* @comment: EurekaService entity
**/
type EurekaService struct {
	eurekaBase *EurekaServer
	logger     *golog.Logger
}

/**
* @author: wing
* @time: 2020/9/4 13:10
* @param:
* @return:
* @comment: Applications entity
**/
type Applications struct {
	VersionsDelta int           `xml:"versions__delta"`
	AppsHashcode  string        `xml:"apps__hashcode"`
	Applications  []Application `xml:"application"`
}

/**
* @author: wing
* @time: 2020/9/4 13:11
* @param:
* @return:
* @comment: Application entity
**/
type Application struct {
	Name      string     `xml:"name"`
	Instances []Instance `xml:"instance"`
}

/**
* @author: wing
* @time: 2020/9/4 13:11
* @param:
* @return:
* @comment: Instance entity
**/
type Instance struct {
	InstanceId                    string          `xml:"instanceId" json:"instanceId"`
	HostName                      string          `xml:"hostName" json:"hostName"`
	App                           string          `xml:"app" json:"app"`
	IpAddr                        string          `xml:"ipAddr" json:"ipAddr"`
	Status                        string          `xml:"status" json:"status"`
	OverriddenStatus              string          `xml:"overriddenstatus,omitempty" json:"overriddenstatus,omitempty"`
	Port                          *Port           `xml:"port,omitempty" json:"port,omitempty"`
	SecurePort                    *Port           `xml:"securePort,omitempty" json:"securePort,omitempty"`
	CountryId                     int             `xml:"countryId,omitempty" json:"countryId,omitempty"`
	DataCenterInfo                *DataCenterInfo `xml:"dataCenterInfo" json:"dataCenterInfo"`
	LeaseInfo                     *LeaseInfo      `xml:"leaseInfo,omitempty" json:"leaseInfo,omitempty"`
	Metadata                      *MetaData       `xml:"metadata,omitempty" json:"metadata,omitempty"`
	HomePageUrl                   string          `xml:"homePageUrl,omitempty" json:"homePageUrl,omitempty"`
	StatusPageUrl                 string          `xml:"statusPageUrl" json:"statusPageUrl"`
	HealthCheckUrl                string          `xml:"healthCheckUrl,omitempty" json:"healthCheckUrl,omitempty"`
	VipAddress                    string          `xml:"vipAddress" json:"vipAddress"`
	SecureVipAddress              string          `xml:"secureVipAddress,omitempty" json:"secureVipAddress,omitempty"`
	IsCoordinatingDiscoveryServer bool            `xml:"isCoordinatingDiscoveryServer,omitempty" json:"isCoordinatingDiscoveryServer,omitempty"`
	LastUpdatedTimestamp          int             `xml:"lastUpdatedTimestamp,omitempty" json:"lastUpdatedTimestamp,omitempty"`
	LastDirtyTimestamp            int             `xml:"lastDirtyTimestamp,omitempty" json:"lastDirtyTimestamp,omitempty"`
	ActionType                    string          `xml:"actionType,omitempty" json:"actionType,omitempty"`
}

/**
* @author: wing
* @time: 2020/9/4 21:32
* @param:
* @return:
* @comment: InstanceObject entity
**/
type InstanceObject struct {
	Instance *Instance `xml:"instance" json:"instance"`
}

/**
* @author: wing
* @time: 2020/9/4 13:11
* @param:
* @return:
* @comment: Port entity
**/
type Port struct {
	Port    int  `xml:",chardata" json:"$"`
	Enabled bool `xml:"enabled,attr" json:"@enabled"`
}

/**
* @author: wing
* @time: 2020/9/4 13:12
* @param:
* @return:
* @comment: DataCenterInfo entity
**/
type DataCenterInfo struct {
	Name     string              `xml:"name" json:"name"`
	Class    string              `xml:"class,attr" json:"@class"`
	Metadata *DataCenterMetadata `xml:"metadata,omitempty" json:"metadata,omitempty"`
}

/**
* @author: wing
* @time: 2020/9/4 13:12
* @param:
* @return:
* @comment: LeaseInfo entity
**/
type LeaseInfo struct {
	RenewalIntervalInSecs  int  `xml:"renewalIntervalInSecs,omitempty" json:"renewalIntervalInSecs,omitempty"`
	DurationInSecs         int  `xml:"durationInSecs,omitempty" json:"durationInSecs,omitempty"`
	RegistrationTimestamp  int  `xml:"registrationTimestamp,omitempty" json:"registrationTimestamp,omitempty"`
	LastRenewalTimestamp   int  `xml:"lastRenewalTimestamp,omitempty" json:"lastRenewalTimestamp,omitempty"`
	EvictionTimestamp      int  `xml:"evictionTimestamp,omitempty" json:"evictionTimestamp,omitempty"`
	ServiceUpTimestamp     int  `xml:"serviceUpTimestamp,omitempty" json:"serviceUpTimestamp,omitempty"`
	EvictionDurationInSecs uint `xml:"evictionDurationInSecs,omitempty" json:"evictionDurationInSecs,omitempty"`
}

/**
* @author: wing
* @time: 2020/9/4 13:12
* @param:
* @return:
* @comment: MetaData entity
**/
type MetaData struct {
	Map   map[string]string
	Class string
}

/**
* @author: wing
* @time: 2020/9/4 13:13
* @param:
* @return:
* @comment: DataCenterMetadata entity
**/
type DataCenterMetadata struct {
	AmiLaunchIndex   string `xml:"ami-launch-index,omitempty" json:"ami-launch-index,omitempty"`
	LocalHostname    string `xml:"local-hostname,omitempty" json:"local-hostname,omitempty"`
	AvailabilityZone string `xml:"availability-zone,omitempty" json:"availability-zone,omitempty"`
	InstanceId       string `xml:"instance-id,omitempty" json:"instance-id,omitempty"`
	PublicIpv4       string `xml:"public-ipv4,omitempty" json:"public-ipv4,omitempty"`
	PublicHostname   string `xml:"public-hostname,omitempty" json:"public-hostname,omitempty"`
	AmiManifestPath  string `xml:"ami-manifest-path,omitempty" json:"ami-manifest-path,omitempty"`
	LocalIpv4        string `xml:"local-ipv4,omitempty" json:"local-ipv4,omitempty"`
	Hostname         string `xml:"hostname,omitempty" json:"hostname,omitempty"`
	AmiId            string `xml:"ami-id,omitempty" json:"ami-id,omitempty"`
	InstanceType     string `xml:"instance-type,omitempty" json:"instance-type,omitempty"`
}

/**
* @author: wing
* @time: 2020/9/4 13:13
* @param:
* @return:
* @comment: create 1 eureka service
**/
func NewEurekaService(eurekaBase *EurekaServer) *EurekaService {
	return &EurekaService{
		eurekaBase: eurekaBase,
		logger:     config.Global.Logger,
	}
}

/**
* @author: wing
* @time: 2020/9/4 13:10
* @param:
* @return:
* @comment: special logger
**/
func (s *EurekaService) SetLogger(logger *golog.Logger) {
	s.logger = logger
}

/**
* @author: wing
* @time: 2020/9/4 13:14
* @param:
* @return:
* @comment: local service register to eureka
**/
func (s *EurekaService) Register(appName string, instance *Instance) error {
	ep := "/apps" + "/" + appName
	contentType := "application/json"
	instObject := &InstanceObject{
		Instance: instance,
	}
	instanceByte, err := json.Marshal(instObject)
	if err != nil {
		s.logger.Errorf("EurekaService.Register.1:%s", err)
		return err
	}
	if _, err := s.eurekaBase.Post(ep, contentType, bytes.NewReader(instanceByte)); err != nil {
		s.logger.Errorf("EurekaService.Register.2:%s", err)
		return err
	}
	return nil
}

/**
* @author: wing
* @time: 2020/9/4 13:14
* @param:
* @return:
* @comment: local unregister register to eureka
**/
func (s *EurekaService) UnRegister(appName string, instanceId string) error {
	ep := "/apps" + "/" + appName + "/" + instanceId
	_, err := s.eurekaBase.Delete(ep, nil)
	if err != nil {
		s.logger.Errorf("EurekaService.UnRegister:%s", err)
		return err
	}
	return nil
}

/**
* @author: wing
* @time: 2020/9/4 13:14
* @param:
* @return:
* @comment: pause service in eureka
**/
func (s *EurekaService) PauseService(appName string, instanceId string) error {
	ep := "/apps" + "/" + appName + "/" + instanceId + "status?value=OUT_OF_SERVICE"
	_, err := s.eurekaBase.Put(ep, nil)
	if err != nil {
		s.logger.Errorf("EurekaService.PauseService:%s", err)
		return err
	}
	if s.eurekaBase.HttpResponse.StatusCode == http.StatusInternalServerError {
		return errors.New("Pause service: appName.instanceId failed! ")
	}
	return nil
}

/**
* @author: wing
* @time: 2020/9/4 13:14
* @param:
* @return:
* @comment: resume service in eureka
**/
func (s *EurekaService) ResumeService(appName string, instanceId string) error {
	ep := "/apps" + "/" + appName + "/" + instanceId + "status?value=UP"
	_, err := s.eurekaBase.Delete(ep, nil)
	if err != nil {
		s.logger.Errorf("EurekaService.ResumeService:%s", err)
		return err
	}
	if s.eurekaBase.HttpResponse.StatusCode == http.StatusInternalServerError {
		return errors.New("Resume service: appName.instanceId failed! ")
	}
	return nil
}

/**
* @author: wing
* @time: 2020/9/4 9:12
* @param:
* @return:
* @comment: send heartbeat to eureka
**/
func (s *EurekaService) SendHeartbeat(appName string, instance *Instance, interval int, done chan bool) {
	ep := "/apps" + "/" + appName + "/" + instance.InstanceId + "?status=UP"
	go func() {
		for {
			select {
			case <-done:
				s.logger.Warnf("Stop EurekaService.SendHeartbeat! ")
				break
			default:

			}
			_, err := s.eurekaBase.Put(ep, nil)
			if err != nil || s.eurekaBase.HttpResponse.StatusCode != http.StatusOK {
				if err := s.Register(appName, instance); err != nil {
					s.logger.Errorf("EurekaService.SendHeartbeat:%s", err)
				}
			}
			s.logger.Debugf("EurekaService.SendHeartbeat: %s", time.Now().Format("2006-01-02 15:04:05"))
			time.Sleep(time.Duration(interval) * time.Second)
		}
	}()
}

/**
* @author: wing
* @time: 2020/9/4 9:12
* @param:
* @return: Application
* @comment: get all application from eureka
**/
func (s *EurekaService) GetApplications() *Applications {
	ep := "/apps"
	resByte, err := s.eurekaBase.Get(ep)
	if err != nil {
		s.logger.Errorf("EurekaService.GetApplications.1:%s", err)
		return nil
	}
	applications := &Applications{}
	if err := xml.Unmarshal(resByte, applications); err != nil {
		s.logger.Errorf("EurekaService.GetApplications.2:%s", err)
		return nil
	}
	return applications
}

/**
* @author: wing
* @time: 2020/9/4 13:15
* @param:
* @return:
* @comment: get application from eureka
**/
func (s *EurekaService) GetApplication(appName string) *Application {
	ep := "/apps" + "/" + appName
	resByte, err := s.eurekaBase.Get(ep)
	if err != nil {
		s.logger.Errorf("EurekaService.GetApplication.1:%s", err)
		return nil
	}
	application := &Application{}
	if err := xml.Unmarshal(resByte, application); err != nil {
		s.logger.Errorf("EurekaService.GetApplication.2:%s", err)
		return nil
	}
	return application
}

/**
* @author: wing
* @time: 2020/9/4 13:15
* @param:
* @return:
* @comment: get instance from eureka
**/
func (s *EurekaService) GetInstance(appName string, instanceId string) *Instance {
	ep := "/apps" + "/" + appName + "/" + instanceId
	resByte, err := s.eurekaBase.Get(ep)
	if err != nil {
		s.logger.Errorf("EurekaService.GetInstance.1:%s", err)
		return nil
	}
	instance := &Instance{}
	if err := xml.Unmarshal(resByte, instance); err != nil {
		s.logger.Errorf("EurekaService.GetInstance.2:%s", err)
		return nil
	}
	return instance
}

/**
* @author: wing
* @time: 2020/9/4 13:16
* @param:
* @return:
* @comment: get instance from eureka by instanceid
**/
func (s *EurekaService) GetInstanceById(instanceId string) *Instance {
	ep := "/instances" + "/" + instanceId
	resByte, err := s.eurekaBase.Get(ep)
	if err != nil {
		s.logger.Errorf("EurekaService.GetInstanceById.1:%s", err)
		return nil
	}
	instance := &Instance{}
	if err := xml.Unmarshal(resByte, instance); err != nil {
		s.logger.Errorf("EurekaService.GetInstanceById.2:%s", err)
		return nil
	}
	return instance
}

/**
* @author: wing
* @time: 2020/9/4 13:16
* @param:
* @return:
* @comment: get instance from eureka by vip
**/
func (s *EurekaService) GetInstanceByVip(vip string) *Applications {
	ep := "/vips" + "/" + vip
	resByte, err := s.eurekaBase.Get(ep)
	if err != nil {
		s.logger.Errorf("EurekaService.GetInstanceByVip.1:%s", err)
		return nil
	}
	applications := &Applications{}
	if err := xml.Unmarshal(resByte, applications); err != nil {
		s.logger.Errorf("EurekaService.GetInstanceByVip.2:%s", err)
		return nil
	}
	return applications
}

/**
* @author: wing
* @time: 2020/9/4 13:16
* @param:
* @return:
* @comment: get instance from eureka by svip
**/
func (s *EurekaService) GetInstanceBySvip(svip string) *Applications {
	ep := "/svips" + "/" + svip
	resByte, err := s.eurekaBase.Get(ep)
	if err != nil {
		s.logger.Errorf("EurekaService.GetInstanceBySvip.1:%s", err)
		return nil
	}
	applications := &Applications{}
	if err := xml.Unmarshal(resByte, applications); err != nil {
		s.logger.Errorf("EurekaService.GetInstanceBySvip.2:%s", err)
		return nil
	}
	return applications
}

/**
* @author: wing
* @time: 2020/9/4 13:17
* @param:
* @return:
* @comment: create 1 remote service client
**/
func (s *EurekaService) GetFeignClient(appName string, contextPath string) *FeignClient {
	return NewFeignClient(s, appName, contextPath)
}
