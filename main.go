package main

/**
  @author: wing
  @date: 2020/9/4
  @comment:
**/
import (
	"go-eureka/client"
	"go-eureka/config"
	"os"
)

func main() {
	log := config.Global.Logger
	signal := make(chan os.Signal)
	exit := make(chan bool)
	c := client.GetServer(exit)
	//webserver 协程
	go func() {
		log.Infof("service start ...")
		c.Serving()
	}()
	//主进程
	go func() {
		<-signal
		log.Infof("service stop ...")
		exit <- true
	}()
	<-exit
	close(exit)
	log.Infof("exit main process!")
}
