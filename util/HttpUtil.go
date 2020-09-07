package util

/**
  @author: wing
  @date: 2020/9/4
  @comment:
**/
import (
	"strings"
)

/**
* @author: wing
* @time: 2020/9/4 9:50
* @param:
* @return:
* @comment: endpoint must be start with '/',and not end with '/',this function may fix it
**/
func FixEndpoint(endpoint string) string {
	var res string
	rep := []rune(endpoint)
	if strings.HasSuffix(endpoint, "/") {
		rep = rep[:len(rep)-1]
	}
	if !strings.HasPrefix(endpoint, "/") {
		res += "/"
	}
	res += string(rep)
	return res
}
