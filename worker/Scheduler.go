package worker

import "corntab/common"

type scheduler struct {
	jobEventChan chan *common.JobEvent
}

var (
	S_scheduler *scheduler
)

func (s *scheduler) scheduleLoop()  {
	
}

func InitScheduler() (err error) {
	s := &scheduler{
		jobEventChan: make(chan *common.JobEvent, 1000),
	}
	S_scheduler = s
	//启动协程开始调度循环
	go S_scheduler.scheduleLoop()
	return
}