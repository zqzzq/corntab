package main

import (
	"go.etcd.io/etcd/clientv3"
	"time"
	"fmt"
)

func main()  {

	//etcd客户端配置
	etcd_conf := clientv3.Config{
		Endpoints: []string{"etcd.10.110.25.114.xip.io:81"},
		DialTimeout: 5 * time.Second,
	}

	//建立连接
	client, err := clientv3.New(etcd_conf)
	if err != nil{
		fmt.Println(err)
		return
	}

	client = client




}