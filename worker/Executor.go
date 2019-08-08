package worker

import (
	"corntab/common"
	"os/exec"
	"context"
	"time"
)

type Executor struct {

}

var (
	S_executor *Executor
)

func (e *Executor)ExecuteJob(info *common.JobExecInfo) {
	go func() {
		cmd := exec.CommandContext(context.TODO(), "/bin/bash", "-c", info.Job.Command)
		startTime := time.Now()
		output, err := cmd.CombinedOutput()
		endTime := time.Now()
		jer := common.BuildJobExecResult(info, string(output), err, startTime, endTime)
		S_scheduler.jobExecResultChan <- jer
	}()

}

func InitExecutor()  {
	S_executor = &Executor{}
}
