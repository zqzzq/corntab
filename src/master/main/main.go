package main

import (
	"runtime"
	"corntab/src/master"
	"fmt"
	"flag"
)

var confPath string

func initArgs()  {
	//master -config ./master.json
	flag.StringVar(&confPath, "config", "src/master/main/master.json", "指定master.json")
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
	if err := master.InitConfig(confPath);err != nil{
		fmt.Println(err)
		return
	}
	//初始化任务管理器
	if err := master.InitJobMgr();err != nil{
		fmt.Println(err)
		return
	}

	//启动HTTP服务
	if err := master.InitApiServer();err != nil{
		fmt.Println(err)
		return
	}



	
}