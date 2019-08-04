package master

import (
	"go.etcd.io/etcd/clientv3"
	"time"
	"corntab/common"
	"encoding/json"
	"context"
)

type JobMgr struct {
	client clientv3.Client
	Kv clientv3.KV
	Lease clientv3.Lease
}

var (
	S_jobMgr *JobMgr
)
//初始化任务管理器 建立与etcd的连接
func InitJobMgr() (err error) {
	//etcd客户端配置
	etcd_conf := clientv3.Config{
		Endpoints: []string{S_config.EtcdEndpoints},
		DialTimeout: time.Duration(S_config.EtcdDialTimeout) * time.Millisecond,
	}

	//建立连接
	cli, err := clientv3.New(etcd_conf)
	if err != nil{
		return
	}
	S_jobMgr = &JobMgr{
		client: *cli,
		Kv: clientv3.NewKV(cli),
		Lease: clientv3.NewLease(cli),
	}
	return
}

func (jobMgr *JobMgr)SaveJob(job *common.Job) (oldJob *common.Job,err error) {
	//把任务保存到/cron/jobs/任务名 -> json

	jobKey := "/cron/jobs/" + job.Name
	jobVal,err := json.Marshal(job)
	if err != nil{
		return
	}
	//保存到etcd中
	putResp, err := jobMgr.Kv.Put(context.TODO(), jobKey, string(jobVal), clientv3.WithPrevKV())
	if err != nil{
		return
	}
	//更新操作 有旧值
	if putResp.PrevKv != nil{
		json.Unmarshal(putResp.PrevKv.Value, oldJob)
		return
	}
	return
}
