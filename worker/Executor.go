package worker

import (
	"corntab/common"
	"os/exec"
	"time"
)

type Executor struct {

}

var (
	S_executor *Executor
)

func (e *Executor)ExecuteJob(info *common.JobExecInfo) {
	go func() {
		//获取分布式锁
		var (
			startTime time.Time
			endTime time.Time
			output []byte
			err error
		)
		jobLock := S_jobMgr.getJobLock(info.Job.Name)
		if err = jobLock.tryLock();err == nil{//上锁成功
			cmd := exec.CommandContext(info.Ctx, "/bin/bash", "-c", info.Job.Command)
			startTime = time.Now()
			output, err = cmd.CombinedOutput()
			endTime = time.Now()
		}
		jer := common.BuildJobExecResult(info, string(output), err, startTime, endTime)
		S_scheduler.jobExecResultChan <- jer
		jobLock.unLock()
	}()
}

func InitExecutor()  {
	S_executor = &Executor{}
}
