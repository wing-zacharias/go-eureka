package config

/**
  @author: wing
  @date: 2020/9/4
  @comment:
**/
import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kataras/golog"
)

var Global *Conf

type Conf struct {
	Logger       *golog.Logger
	GlobalConfig *GlobalConf `json:"global"`
	EurekaConfig *EurekaConf `json:"eureka"`
}

/**
* @author: wing
* @time: 2020/9/4 12:37
* @param:
* @return:
* @comment: autodetect entity
**/
type AutoDetect struct {
	Allow       bool   `json:"allow"`
	Network     string `json:"network"`
	Port        string `json:"port"`
	ContextPath string `json:"contextPath"`
	Timeout     int    `json:"timeout"`
}

/**
* @author: wing
* @time: 2020/9/4 12:38
* @param:
* @return:
* @comment: eureka config entity
**/
type EurekaConf struct {
	AutoDetect  *AutoDetect `json:"autoDetect"`
	EurekaNodes []string    `json:"eurekaNodes"`
}

/**
* @author: wing
* @time: 2020/9/4 15:35
* @param:
* @return:
* @comment: global config entity
**/
type GlobalConf struct {
	Common *CommonConf `json:"common"`
	Log    *LogConf    `json:"log"`
}

/**
* @author: wing
* @time: 2020/9/4 15:36
* @param:
* @return:
* @comment: common entity
**/
type CommonConf struct {
	Listen       string `json:"listen"`
	ContextPath  string `json:"contextPath"`
	Port         string `json:"port"`
	AppName      string `json:"appName"`
	AutoRegister bool   `json:"autoRegister"`
}

/**
* @author: wing
* @time: 2020/9/4 15:36
* @param:
* @return:
* @comment: log entity
**/
type LogConf struct {
	Console *ConsoleConf `json:"console"`
	LogFile *LogFileConf `json:"logFile"`
}

/**
* @author: wing
* @time: 2020/9/4 15:36
* @param:
* @return:
* @comment: console entity
**/
type ConsoleConf struct {
	LogLevel string `json:"level"`
}

/**
* @author: wing
* @time: 2020/9/4 15:36
* @param:
* @return:
* @comment: log file entity
**/
type LogFileConf struct {
	Enable   bool   `json:"enable"`
	Path     string `json:"path"`
	LogLevel string `json:"level"`
	MaxLine  int    `json:"maxLine"`
	MaxDay   int    `json:"maxDay"`
}

/**
* @author: wing
* @time: 2020/9/4 15:36
* @param:
* @return:
* @comment: reload config
**/
func (g *Conf) Reload() {
	currentFilePath, _ := os.Executable()
	configFile := filepath.Join(filepath.Dir(currentFilePath), "/config/config.json")
	conf, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(conf, &Global); err != nil {
		panic(err)
	}
}

/**
* @author: wing
* @time: 2020/9/4 15:37
* @param:
* @return:
* @comment: auto read config
**/
func init() {
	logFile := &LogFileConf{}
	console := &ConsoleConf{}
	log := &LogConf{
		Console: console,
		LogFile: logFile,
	}
	common := &CommonConf{}
	globalConfig := &GlobalConf{
		Common: common,
		Log:    log,
	}
	autoDetect := &AutoDetect{}
	eurekaConfig := &EurekaConf{
		EurekaNodes: []string{},
		AutoDetect:  autoDetect,
	}
	logger := golog.New()
	Global = &Conf{
		Logger:       logger,
		GlobalConfig: globalConfig,
		EurekaConfig: eurekaConfig,
	}
	Global.Reload()
	Global.Logger.SetLevel(Global.GlobalConfig.Log.Console.LogLevel)
}
