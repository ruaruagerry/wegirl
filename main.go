/*
 * @Author: lingguohua
 * @Date: 2019-08-12 19:46:21
 * @Description:
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
	"wegirl/gamecfg"
	"wegirl/gfunc"
	"wegirl/server"
	"wegirl/servercfg"

	_ "net/http/pprof"

	_ "wegirl/handles/auth"

	log "github.com/sirupsen/logrus"
)

var (
	cfgFilepath = ""
	genInPath   = ""
	genOutPath  = ""
)

var (
	// 显示版本
	showVer = false
	// BuildVersion 编译版本
	BuildVersion string
	//BuildTime 编译时间
	BuildTime string
	//CommitID 提供ID
	CommitID string
)

func init() {
	flag.StringVar(&cfgFilepath, "c", "servercfg/x.json", "specify the config file path name")
	flag.StringVar(&genInPath, "gi", "", "input path")
	flag.StringVar(&genOutPath, "go", "", "output path")
}

func main() {
	// only one thread
	runtime.GOMAXPROCS(1)

	flag.BoolVar(&showVer, "v", false, "show version")

	flag.Parse()

	if showVer {
		fmt.Println("Build Version:", BuildVersion)
		fmt.Println("Build Time:", BuildTime)
		fmt.Println("CommitID:", CommitID)
		os.Exit(0)
	}

	if genInPath != "" && genOutPath != "" {
		gamecfg.Gen(genInPath, genOutPath)
		os.Exit(0)
	}

	if cfgFilepath != "" {
		r := servercfg.ParseConfigFile(cfgFilepath)
		if r != true {
			log.Fatal("can't parse configure file:", cfgFilepath)
		}
	} else {
		log.Fatal("please specify a valid config file path")
	}

	// server startTime
	servercfg.StartTime = int(time.Now().Unix())

	log.Println("try to start  stupid server...")

	// load sensitiveword
	gfunc.LoadSensitiveWordDictionary(servercfg.SensitiveWordFile)

	// start http server
	server.CreateHTTPServer()

	// start cron job
	server.CronJob()

	log.Println("start stupid server ok!")

	if servercfg.Daemon == "yes" {
		waitForSignal()
	} else {
		waitInput()
	}
	return
}

func waitInput() {
	var cmd string
	for {
		_, err := fmt.Scanf("%s\n", &cmd)
		if err != nil {
			//log.Println("Scanf err:", err)
			continue
		}

		switch cmd {
		case "exit", "quit":
			log.Println("exit by user")
			return
		case "gr":
			log.Println("current goroutine count:", runtime.NumGoroutine())
			break
		case "gd":
			pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
			break
		default:
			break
		}
	}
}

func dumpGoRoutinesInfo() {
	log.Println("current goroutine count:", runtime.NumGoroutine())
	// use DEBUG=2, to dump stack like golang dying due to an unrecovered panic.
	pprof.Lookup("goroutine").WriteTo(os.Stdout, 2)
}
