package server

/*
 * @Author: zheng
 * @Date: 2019-08-09 09:22:06
 * @Description:
 */

import (
	"flag"
	"fmt"
	"wegirl/servercfg"
	"wegirl/tables"

	// import mysql
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"
)

var (
	showSQL      = false //是否显示sql
	dbConnect    *xorm.Engine
	logDbConnect *xorm.Engine
)

type dbconnectkey struct {
	Rid     int64 `xorm:"pk autoincr <-"`
	ConnKey string
}

func init() {
	flag.BoolVar(&showSQL, "showsql", false, "显示调用数据库sql")
}

func mysqlStartup() {
	if servercfg.DbIP == "" || servercfg.DbUser == "" || servercfg.DbName == "" || servercfg.DbPassword == "" {
		log.Panic("Must specify the DbServer info in config json")
		return
	}
	log.Infof("connect db addr:%s name:%s account:%s password:%s", servercfg.DbIP, servercfg.DbName, servercfg.DbUser, servercfg.DbPassword)

	conn, err := connect(servercfg.DbIP, servercfg.DbName, servercfg.DbUser, servercfg.DbPassword)
	if err != nil || conn == nil {
		log.Panicf("connect db addr:%s name:%s account:%s password:%s err:%v", servercfg.DbIP, servercfg.DbName, servercfg.DbUser, servercfg.DbPassword, err)
	}

	log.Infof("connect log db addr:%s name:%s account:%s password:%s", servercfg.LogDbIP, servercfg.LogDbName, servercfg.LogDbUser, servercfg.LogDbPassword)

	logConn, err := connect(servercfg.LogDbIP, servercfg.LogDbName, servercfg.LogDbUser, servercfg.LogDbPassword)
	if err != nil || logConn == nil {
		log.Panicf("connect log db addr:%s name:%s account:%s password:%s err:%v", servercfg.LogDbIP, servercfg.LogDbName, servercfg.LogDbUser, servercfg.LogDbPassword, err)
	}

	dbConnect = conn
	logDbConnect = logConn

	if err = createTables(dbConnect); err != nil {
		log.Panicf("CreateTables err:%v", err)
	}
	if err = createLogTables(logDbConnect); err != nil {
		log.Panicf("createLogTables err:%v", err)
	}

	result := checkMysqlKey(dbConnect, fmt.Sprintf("%d", servercfg.ServerID))
	if !result {
		log.Panic("check mysql key failed")
	}
}

// 连接数据库
func connect(addr string, name string, account string, password string) (*xorm.Engine, error) {
	sqlInfo := fmt.Sprintf("%s:%s@tcp(%s)/%s", account, password, addr, name)
	engine, err := xorm.NewEngine("mysql", sqlInfo)

	if err != nil {
		return nil, err
	}

	if showSQL {
		//显示调用sql
		engine.ShowSQL()
	}

	err = engine.Ping()
	if err != nil {
		return nil, err
	}

	return engine, nil
}

// createTables 创建业务表
func createTables(engine *xorm.Engine) error {
	account := &tables.Account{}

	// 创建表
	if err := engine.CreateTables(account); err != nil {
		log.Panicf("CreateTable Player err:%v", err)
		return err
	}

	// 同步表结构
	if err := engine.Sync2(account); err != nil {
		log.Panicf("Syn2 Tables err:%v", err)
		return err
	}

	return nil
}

// createTables 创建运营统计表
func createLogTables(engine *xorm.Engine) error {

	account := &tables.Account{}

	// 创建表
	if err := engine.CreateTables(
		account); err != nil {
		log.Panicf("CreateTable Player err:%v", err)
		return err
	}

	// 同步表结构
	if err := engine.Sync2(
		account); err != nil {
		log.Panicf("Syn2 Tables err:%v", err)
		return err
	}

	return nil
}

// checkMysqlKey 检查数据库的key
func checkMysqlKey(db *xorm.Engine, key string) bool {
	row := &dbconnectkey{}
	if err := db.CreateTables(row); err != nil {
		return false
	}

	has, err := db.Where("rid = 1").Exist(row)
	if err != nil {
		return false
	}

	if has {
		row := &dbconnectkey{}
		row.Rid = 1
		_, err := db.Get(row)
		if err != nil {
			return false
		}

		if row.ConnKey != key {
			log.Printf("row connkey:%s key:%s", row.ConnKey, key)
			return false
		}
	} else {
		row := &dbconnectkey{}
		row.ConnKey = key
		_, err := db.Insert(row)
		if err != nil {
			return false
		}
	}

	return true
}
