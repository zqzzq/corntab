package common

import (
	"encoding/json"
	"strings"
	"github.com/gorhill/cronexpr"
	"time"
)

type Job struct {
	Name string `json:"name"`
	Command string `json:"command"`
	CronExpr string `json:"cronExpr"`
}

type Response struct {
	RespNum int `json:"respNum"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

type JobEvent struct {
	EventType int `json:"eventType"` //delete:1 put:0
	JobInfo *Job `json:"jobInfo"`
}

type JobSchedulePlan struct {
	Job *Job `json:"job"`
	Expr *cronexpr.Expression `json:"expr"`
	NextTime time.Time `json:"nextTime"`
}

type JobExecInfo struct {
	Job *Job
	ScheduleTime time.Time
	ExecTime time.Time
}


func BuildResponse(respNum int, msg string, data interface{}) (resp []byte, err error) {
	r := Response{
		RespNum: respNum,
		Msg: msg,
		Data: data,
	}
	resp, err = json.Marshal(r)
	return
}

func UnPackJob(b []byte) (job *Job, err error) {
	job = &Job{}
	if err = json.Unmarshal(b, job);err !=nil{
		return
	}
	return
}
//从etcd的key中取得任务名
func GetJobName(key string) string {
	return strings.TrimPrefix(key, JOB_SAVE_DIR)
}

func BuildJobEvent(etype int, jinfo *Job) *JobEvent {
	return &JobEvent{
		EventType: etype,
		JobInfo: jinfo,
	}
}

func BuildJobExecInfo(jsp *JobSchedulePlan) *JobExecInfo {
	return &JobExecInfo{
		Job: jsp.Job,
		ScheduleTime: jsp.NextTime,
		ExecTime: time.Now(),
	}
}

func BuildJobSchedulePlan(j *Job) (jsp *JobSchedulePlan, err error) {
	expr, err := cronexpr.Parse(j.CronExpr)
	if err != nil {
		return
	}

	return &JobSchedulePlan{
		Job: j,
		Expr: expr,
		NextTime: expr.Next(time.Now()),
	}, nil
}