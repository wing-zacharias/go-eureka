package util

/**
  @author: wing
  @date: 2020/9/4
  @comment:
**/
import (
	"net"
	"time"
)

/**
* @author: wing
* @time: 2020/9/4 10:27
* @param:
* @return:
* @comment: tcp port test
**/
func Dail(dailIp string, dailPort string, timeout int) (bool, error) {
	conn, err := net.DialTimeout("tcp", dailIp+":"+dailPort, time.Duration(timeout)*time.Millisecond)
	if err != nil {
		return false, err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}()
	return true, nil
}
