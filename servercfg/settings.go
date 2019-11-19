/*
 * @Author: lingguohua
 * @Date: 2019-08-09 17:02:07
 * @Description:
 */

// Package servercfg 游戏服务器运行配置
package servercfg

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/DisposaBoy/JsonConfigReader"
)

// make a copy of this file, rename to settings.go
// then set the correct value for these follow variables
var (
	monitorEstablished = false
	ServerPort         = 3001
	// LogFile              = ""
	Daemon      = "yes"
	StartTime   = 0
	RedisServer = ":6379"
	ServerID    = 0
	URL         = ""

	GameCfgsDir = ""
	DbIP        = "127.0.0.1:3306"
	DbUser      = "user"
	DbPassword  = "123456"
	DbName      = ""

	LogDbIP       = "127.0.0.1:3306"
	LogDbUser     = "user"
	LogDbPassword = "123456"
	LogDbName     = ""

	MQIP       = ":5672"
	MQAccount  = "guest"
	MQPassword = "guest"

	ForTestOnly = false

	SensitiveWordFile = ""
)

var (
	loadedCfgFilePath = ""
)

// ReLoadConfigFile 重新加载配置
func ReLoadConfigFile() bool {
	log.Println("ReLoadConfigFile-----------From File--------:", loadedCfgFilePath)
	if !ParseConfigFile(loadedCfgFilePath) {
		log.Println("ReLoadConfigFile-------------------FAILED")
		return false
	}

	log.Println("ReLoadConfigFile-------------------OK")
	return true
}

// ParseConfigFile 解析配置
func ParseConfigFile(filepath string) bool {
	type DbServerInfo struct {
		DbName     string `json:"dbName"`
		DbIP       string `json:"dbIP"`
		DbUser     string `json:"dbUser"`
		DbPassword string `json:"dbPassword"`
	}

	type Params struct {
		ServerPort int `json:"port"`
		// LogFile           string `json:"log_file"`
		Daemon            string       `json:"daemon"`
		RedisServer       string       `json:"redis_server"`
		ServreID          int          `json:"guid"`
		URL               string       `json:"url"`
		Db                DbServerInfo `json:"db"`
		Dblog             DbServerInfo `json:"db_log"`
		MQIP              string       `json:"mqIP"`
		MQAccount         string       `json:"mqAccount"`
		MQPassword        string       `json:"mqPassword"`
		GameCfgsDir       string       `json:"gameCfgsDir"`
		ForTestOnly       bool         `json:"forTestOnly"`
		SensitiveWordFile string       `json:"SensitiveWordFile"`
	}

	loadedCfgFilePath = filepath

	var params = &Params{}

	f, err := os.Open(filepath)
	if err != nil {
		log.Println("failed to open config file:", filepath)
		return false
	}

	// wrap our reader before passing it to the json decoder
	r := JsonConfigReader.New(f)
	err = json.NewDecoder(r).Decode(params)

	if err != nil {
		log.Println("json un-marshal error:", err)
		return false
	}

	log.Println("-------------------Configure params are:-------------------")
	log.Printf("%+v\n", params)

	// if params.LogFile != "" {
	// 	LogFile = params.LogFile
	// }

	if params.Daemon != "" {
		Daemon = params.Daemon
	}

	if params.ServerPort != 0 {
		ServerPort = params.ServerPort
	}

	if params.RedisServer != "" {
		RedisServer = params.RedisServer
	}

	if params.ServreID != 0 {
		ServerID = params.ServreID
	}

	if params.URL != "" {
		URL = params.URL
	}

	if params.GameCfgsDir != "" {
		GameCfgsDir = params.GameCfgsDir
	}

	if params.MQIP != "" {
		MQIP = params.MQIP
	}

	if params.MQAccount != "" {
		MQAccount = params.MQAccount
	}

	if params.MQPassword != "" {
		MQPassword = params.MQPassword
	}

	DbIP = params.Db.DbIP
	DbUser = params.Db.DbUser
	DbPassword = params.Db.DbPassword
	DbName = params.Db.DbName

	LogDbIP = params.Dblog.DbIP
	LogDbUser = params.Dblog.DbUser
	LogDbPassword = params.Dblog.DbPassword
	LogDbName = params.Dblog.DbName

	// 是否测试用途
	ForTestOnly = params.ForTestOnly

	if params.SensitiveWordFile != "" {
		SensitiveWordFile = params.SensitiveWordFile
	}

	if ServerID == 0 {
		log.Println("Server id 'guid' must not be empty!")
		return false
	}

	if RedisServer == "" {
		log.Println("redis server id  must not be empty!")
		return false
	}

	return true
}
