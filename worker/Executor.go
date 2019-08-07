package worker

import "corntab/common"

type Executor struct {

}

var (
	S_executor *Executor
)

func (e *Executor)ExecuteJob(info *common.JobExecInfo)  {

}

func InitExecutor()  {
	S_executor = &Executor{}
}
