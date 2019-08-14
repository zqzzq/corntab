package master

import (
	"go.etcd.io/etcd/clientv3"
	"time"
	"corntab/src/common"
	"encoding/json"
	"context"
	"fmt"
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

func (jobMgr *JobMgr)SaveJob(job *common.Job) (oldJob common.Job,err error) {
	//把任务保存到/cron/jobs/任务名 -> json
	fmt.Println("save job :", job)
	jobKey := common.JOB_SAVE_DIR + job.Name
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
		json.Unmarshal(putResp.PrevKv.Value, &oldJob)
		return
	}
	return
}

func (jobMgr *JobMgr)DeleteJob(jobName string) (oldJob common.Job,err error) {
	jobKey := common.JOB_DEL_DIR + jobName
	delResp, err := jobMgr.Kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV())
	if err != nil{
		return
	}
	//删除操作 有旧值
	if delResp.PrevKvs != nil{
		json.Unmarshal(delResp.PrevKvs[0].Value, &oldJob)
		return
	}
	return
}

func (jobMgr *JobMgr)ListJob() (jobList []common.Job,err error) {
	jobKey := common.JOB_LIST_DIR
	listResp, err := jobMgr.Kv.Get(context.TODO(), jobKey, clientv3.WithPrefix())
	if err != nil{
		return
	}
	jobList = make([]common.Job, 0)
	if len(listResp.Kvs) != 0 {
		for _, v := range listResp.Kvs{
			j := common.Job{}
			if err := json.Unmarshal(v.Value, &j);err != nil{
				fmt.Println(err)
				continue
			}
			jobList = append(jobList, j)
		}
		return
	}
	return
}

func (jobMgr *JobMgr)KillJob(job string) (oldJob common.Job,err error) {
	//把要杀死的任务保存到/cron/kill/任务名 -> ""
	fmt.Println("kill job :", job)
	jobKey := common.JOB_KILL_DIR + job

	//设置租约
	lresp, err := jobMgr.Lease.Grant(context.TODO(), 1)
	if err != nil{
		return
	}
	//保存到etcd中
	_, err = jobMgr.Kv.Put(context.TODO(), jobKey, "",clientv3.WithLease(lresp.ID))
	if err != nil{
		return
	}
	return
}