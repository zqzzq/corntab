package worker

import (
	"go.etcd.io/etcd/clientv3"
	"context"
	"corntab/src/common"
)

type JobLock struct {
	kv clientv3.KV
	lease clientv3.Lease
	jobName string
	cancel context.CancelFunc
	leaseID clientv3.LeaseID
	isLocked bool
}

func (jobLock *JobLock) tryLock() (err error) {

	//创建租约
	ctx, cancel := context.WithCancel(context.TODO())
	leaseResp, err := jobLock.lease.Grant(ctx, 5)
	if err != nil{
		return
	}
	//续租
	_, err = jobLock.lease.KeepAlive(ctx, leaseResp.ID)
	if err != nil{//续租失败
		jobLock.handleLockFail(leaseResp.ID, cancel)
		return
	}
	//创建txn事务
	txn := jobLock.kv.Txn(context.TODO())
	//锁路径
	lockKey := common.JOB_LOCK_DIR + jobLock.jobName
	//事务强锁
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseResp.ID)))
	txnResp, err := txn.Commit()
	if err != nil{
		jobLock.handleLockFail(leaseResp.ID, cancel)
		return
	}
	if txnResp.Succeeded{//抢锁成功
		jobLock.cancel = cancel
		jobLock.leaseID = leaseResp.ID
		jobLock.isLocked = true
		return
	}
	err = common.LOCK_ALREADY_USED
	jobLock.handleLockFail(leaseResp.ID, cancel)
	return
}

func (jobLock *JobLock)unLock()  {
	if jobLock.isLocked{
		jobLock.cancel()
		jobLock.lease.Revoke(context.TODO(), jobLock.leaseID)
	}
}

func (jobLock *JobLock) handleLockFail(id clientv3.LeaseID, cancel context.CancelFunc)  {
	cancel()
	jobLock.lease.Revoke(context.TODO(), id)
}