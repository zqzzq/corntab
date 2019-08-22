package worker

import (
	"corntab/src/common"
	"time"
	"google.golang.org/grpc"
	"fmt"
	"corntab/src/worker/pb"
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
			err error
			conn *grpc.ClientConn
			resp *pb.ExecuteResponse
		)
		jobLock := S_jobMgr.getJobLock(info.Job.Name)
		if err = jobLock.tryLock();err == nil{//上锁成功
			conn, err = grpc.Dial(info.Job.ExecutorAddr ,grpc.WithInsecure())

			defer conn.Close()
			if err == nil {
				c := pb.NewExecuteServiceClient(conn)
				startTime = time.Now()
				resp, err = c.Execute(info.Ctx, &pb.ExecuteRequest{Params: info.Job.Params})
				endTime = time.Now()
				if err != nil {
					fmt.Println("执行失败：", err)
				}
			}else {
				fmt.Println("Can't connect to executor: ", err)
			}
		}
		if resp == nil{
			resp = &pb.ExecuteResponse{Output: ""}
		}
		jer := common.BuildJobExecResult(info, resp.Output, err, startTime, endTime)
		S_scheduler.jobExecResultChan <- jer
		jobLock.unLock()
	}()
}

func InitExecutor()  {
	S_executor = &Executor{}
}
