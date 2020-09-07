package erpc

/**
  @author: wing
  @date: 2020/9/4
  @comment:
**/
import (
	"go-eureka/config"
	"net/url"

	"github.com/kataras/golog"
)

/**
* @author: wing
* @time: 2020/9/4 9:14
* @param:
* @return:
* @comment: Cluster entity: main service provider
**/
type Cluster struct {
	Leader      string   `json:"leader"`
	EurekaNodes []string `json:"eurekaNodes"`
	logger      *golog.Logger
}

/**
* @author: wing
* @time: 2020/9/4 9:15
* @param: eurekaNodes: eureka nodes,at lease 1
* @return: Cluster
* @comment:  create 1 eureka cluster instance
**/
func GetCluster(eurekaNodes []string) *Cluster {
	if len(eurekaNodes) != 0 {
		return &Cluster{
			Leader:      eurekaNodes[0],
			EurekaNodes: eurekaNodes,
			logger:      config.Global.Logger,
		}
	}
	return nil
}

/**
* @author: wing
* @time: 2020/9/4 13:06
* @param:
* @return:
* @comment: special logger
**/
func (c *Cluster) SetLogger(logger *golog.Logger) {
	c.logger = logger
}

/**
* @author: wing
* @time: 2020/9/4 9:18
* @param: num:number of eureka cluster list
* @return:
* @comment: special cluster leader
**/
func (c *Cluster) updateLeaderByOrder(num int) {
	if num > 0 && num < len(c.EurekaNodes) {
		c.Leader = c.EurekaNodes[num]
		c.logger.Warnf("Special cluster leader: %s", c.Leader)
	}
}

/**
* @author: wing
* @time: 2020/9/4 9:20
* @param: eurekaUrl:eureka url
* @return:
* @comment: special cluster leader
**/
func (c *Cluster) updateLeaderByUrl(eurekaUrl string) {
	c.Leader = eurekaUrl
	c.logger.Warnf("Special cluster leader: %s", c.Leader)
}

/**
* @author: wing
* @time: 2020/9/4 9:22
* @param: rUrl:request of url
* @return:
* @comment: update cluster leader from request
**/
func (c *Cluster) updateLeaderByRequestUrl(rUrl *url.URL) {
	var leader string
	if rUrl.Scheme == "" {
		leader = "http://" + rUrl.Host
	} else {
		leader = rUrl.Scheme + "://" + rUrl.Host
	}
	c.updateLeaderByUrl(leader)
	c.logger.Warnf("Special cluster leader: %s", c.Leader)
}
