package worker

import (
	"corntab/common"
	"time"
	"fmt"
)

type Scheduler struct {
	jobEventChan chan *common.JobEvent
	jobPlanTable map[string]*common.JobSchedulePlan//任务调度计划表
	jobExecTable map[string] *common.JobExecInfo//任务执行计划表
}

var (
	S_scheduler *Scheduler
)

func (s *Scheduler) tryExecJob(jsp *common.JobSchedulePlan)  {
	if _, isExec := s.jobExecTable[jsp.Job.Name]; isExec{//正在执行
	fmt.Println("执行尚未完成:", jsp.Job.Name)
		return
	}

	jobExecInfo := common.BuildJobExecInfo(jsp)
	s.jobExecTable[jsp.Job.Name] = jobExecInfo
	fmt.Println("开始执行任务", jobExecInfo.Job.Name, jobExecInfo.ScheduleTime,jobExecInfo.ExecTime)
	S_executor.ExecuteJob(jobExecInfo)

}

//执行到期任务并获取距离最近的下次执行时间
func (s *Scheduler) trySchedule() (scheduleAfter time.Duration) {
	if len(s.jobPlanTable)==0 {
		return time.Second
	}
	var nearScheduleTime *time.Time
	now := time.Now()
	for _, jsp := range s.jobPlanTable{
		if jsp.NextTime.Before(now) || jsp.NextTime.Equal(now){
			S_scheduler.tryExecJob(jsp)
			jsp.NextTime = jsp.Expr.Next(now)
		}
		if nearScheduleTime == nil || jsp.NextTime.Before(*nearScheduleTime){
			nearScheduleTime = &jsp.NextTime
		}
	}
	return nearScheduleTime.Sub(now)
}

func (s *Scheduler) handleJobEvent( e *common.JobEvent ) (err error) {
	switch e.EventType {
	case common.JOB_ENEVT_SAVE:
		jsp, err := common.BuildJobSchedulePlan(e.JobInfo)
		if err != nil{
			return err
		}
		s.jobPlanTable[e.JobInfo.Name] = jsp
	case common.JOB_ENEVT_DELETE:
		if _, isExist := s.jobPlanTable[e.JobInfo.Name]; isExist{
			delete(s.jobPlanTable, e.JobInfo.Name)
		}

	}
	return
}

func (s *Scheduler) scheduleLoop()  {

	scheduleAfter := s.trySchedule()
	delayTime := time.NewTimer(scheduleAfter)
	
	for  {
		jobEvent := &common.JobEvent{}
		select {
		case jobEvent = <- s.jobEventChan:
			if err := s.handleJobEvent(jobEvent);err != nil{
				continue
			}
		case <- delayTime.C:
		}
		scheduleAfter = s.trySchedule()
		delayTime = time.NewTimer(scheduleAfter)
	}
}

func InitScheduler() (err error) {
	s := &Scheduler{
		jobEventChan: make(chan *common.JobEvent, 1000),
		jobPlanTable: make(map[string]*common.JobSchedulePlan),
		jobExecTable: make(map[string]*common.JobExecInfo),
	}
	S_scheduler = s
	//启动协程开始调度循环
	go S_scheduler.scheduleLoop()
	return
}

func (s *Scheduler) PushJobEvent(e *common.JobEvent)  {
	s.jobEventChan <- e
}