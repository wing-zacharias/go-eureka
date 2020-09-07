package network

/**
  @author: wing
  @date: 2020/9/4
  @comment:
**/
import (
	"go-eureka/util"
)

/**
* @author: wing
* @time: 2020/9/4 10:39
* @param:
* @return:
* @comment: port status entity
**/
type PortStatus struct {
	Ip   string
	Port string
	Open bool
}

/**
* @author: wing
* @time: 2020/9/4 10:39
* @param: scanIps: ipv4 list,scanPorts: need detective ports,timeout: nanosecond
* @return: PortStatus
* @comment: scan ports of servers
**/
func Scan(scanIps []string, scanPorts []string, timeout int) []PortStatus {
	var portsStatus []PortStatus
	for _, ip := range scanIps {
		for _, port := range scanPorts {
			b, _ := util.Dail(ip, port, timeout)
			ps := PortStatus{
				Ip:   ip,
				Port: port,
				Open: b,
			}
			portsStatus = append(portsStatus, ps)
		}
	}
	return portsStatus
}

/**
* @author: wing
* @time: 2020/9/4 10:43
* @param: timeout: nanosecond
* @return:
* @comment: scan local network using mask 255.255.255.0
**/
func FastScan(localIp string, port string) []PortStatus {
	return FullScan(localIp, 24, port, 1000)
}

/**
* @author: wing
* @time: 2020/9/4 12:44
* @param: timeout: nanosecond
* @return: all scan result
* @comment: scan local network hosts
**/
func FullScan(localIp string, maskLen int, port string, timeout int) []PortStatus {
	var portsStatus []PortStatus
	ips := util.GetNetIpList(localIp, maskLen)
	for _, ip := range ips {
		b, _ := util.Dail(ip, port, timeout)
		portStatus := PortStatus{
			Ip:   ip,
			Port: port,
			Open: b,
		}
		portsStatus = append(portsStatus, portStatus)
	}
	return portsStatus
}

/**
* @author: wing
* @time: 2020/9/4 10:43
* @param: timeout: nanosecond
* @return: open port hosts
* @comment: scan local network hosts
**/
func PortOpenFullScan(localIp string, maskLen int, port string, timeout int) []PortStatus {
	var portsStatus []PortStatus
	ips := util.GetNetIpList(localIp, maskLen)
	for _, ip := range ips {
		b, _ := util.Dail(ip, port, timeout)
		if b {
			portStatus := PortStatus{
				Ip:   ip,
				Port: port,
				Open: b,
			}
			portsStatus = append(portsStatus, portStatus)
		}
	}
	return portsStatus
}
