package main

import (
	"runtime"
	"fmt"
	"flag"
	"corntab/worker"
	"time"
)

var confPath string

func initArgs()  {
	//worker -config ./worker.json
	flag.StringVar(&confPath, "config", "worker/main/worker.json", "指定worker.json")
	flag.Parsed()
}

func initEnv()  {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main()  {
	//初始化命令行参数
	initArgs()
	//初始化线程
	initEnv()
	//加载配置文件
	if err := worker.InitConfig(confPath);err != nil{
		fmt.Println(err)
		return
	}
	//初始化任务执行器
	worker.InitExecutor()
	//初始化任务调度器
	if err := worker.InitScheduler();err != nil{
		fmt.Println(err)
		return
	}
	//初始化任务管理器
	if err := worker.InitJobMgr();err != nil{
		fmt.Println(err)
		return
	}
	for{
		time.Sleep(1)
	}

}