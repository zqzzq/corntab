package worker

import (
	"go.etcd.io/etcd/clientv3"
	"time"
	"corntab/src/common"
	"context"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"fmt"
)

type JobMgr struct {
	client clientv3.Client
	Kv clientv3.KV
	Lease clientv3.Lease
	watcher clientv3.Watcher
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
		watcher: clientv3.NewWatcher(cli),
	}

	//启动任务监听
	if err = S_jobMgr.WatchJobs();err != nil{
		return
	}
	S_jobMgr.WatchKill()
	return
}
//监视任务变化
func (jobMgr *JobMgr) WatchJobs() (err error) {
	jobKey := common.JOB_SAVE_DIR
	getResp, err := jobMgr.Kv.Get(context.TODO(), jobKey, clientv3.WithPrefix())
	if err != nil{
		return
	}
	//获取到当前的所有任务
	for _, kv := range getResp.Kvs{
		if job, err := common.UnPackJob(kv.Value);err == nil{
			event := common.BuildJobEvent(common.JOB_ENEVT_SAVE, job)
			fmt.Println("推送已有任务事件：", event)
			S_scheduler.PushJobEvent(event)
		}
	}
	go func() {
		//从当前版本开始监听
		watchStartV := getResp.Header.Revision + 1
		wchan := S_jobMgr.watcher.Watch(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithRev(watchStartV), clientv3.WithPrefix())
		for wresp := range wchan{
			for _, wEvent := range wresp.Events{
				event := &common.JobEvent{}
				switch wEvent.Type{
				case mvccpb.PUT:
					if j, err := common.UnPackJob(wEvent.Kv.Value);err == nil{
						event = common.BuildJobEvent(common.JOB_ENEVT_SAVE, j)
					}
				case mvccpb.DELETE:
					event = common.BuildJobEvent(common.JOB_ENEVT_DELETE, &common.Job{Name: common.GetJobName(string(wEvent.Kv.Key))})
				}
				fmt.Println("监听到事件：", event)
				S_scheduler.PushJobEvent(event)
			}
		}
	}()
	return
}

func (jobMgr *JobMgr) WatchKill(){
	go func() {
		//从当前版本开始监听
		wchan := S_jobMgr.watcher.Watch(context.TODO(), common.JOB_KILL_DIR, clientv3.WithPrefix())
		for wresp := range wchan{
			for _, wEvent := range wresp.Events{
				event := &common.JobEvent{}
				switch wEvent.Type{
				case mvccpb.PUT:
					killJobName := common.GetKillName(string(wEvent.Kv.Key))
					event.EventType = common.JOB_ENEVT_KILL
					event.JobInfo = &common.Job{Name: killJobName}
					fmt.Println("监听到杀死事件：", event)
					S_scheduler.PushJobEvent(event)
				case mvccpb.DELETE:
				}
			}
		}
	}()
}

func (jobMgr *JobMgr) getJobLock(jobName string) *JobLock {
	return &JobLock{
		kv: jobMgr.Kv,
		lease: jobMgr.Lease,
		jobName: jobName,
	}
}


