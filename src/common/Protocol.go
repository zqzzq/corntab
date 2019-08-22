package common

import (
	"encoding/json"
	"strings"
	"github.com/gorhill/cronexpr"
	"time"
	"context"
)

type Job struct {
	Name string `json:"name"`
	Params string `json:"params"`
	ExecutorAddr string `json:"executorAddr"`
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
	PlanTime time.Time
	ScheduleTime time.Time
	Ctx context.Context
	Cancel context.CancelFunc
}

type JobExecResult struct {
	Info *JobExecInfo
	Output string
	Err error
	StartTime time.Time
	EndTime time.Time
}

type JobLog struct {
	JobName string `bson:"jobName"`
	Params string `bson:"params"`
	Err string `bson:"err"`
	OutPut string `bson:"outPut"`
	PlanTime time.Time `bson:"planTime"`
	ScheduleTime time.Time `bson:"scheduleTime"`
	StartTime time.Time `bson:"startTime"`
	EndTime time.Time `bson:"endTime"`
}

type JobLogBatch struct {
	JobLogs []interface{}
}

func BuildResponse(respNum int, msg string, data interface{}) (resp []byte, err error) {
	r := Response{
		RespNum: respNum,
		Msg: msg,
		Data: data,
	}
	return json.Marshal(r)

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

func GetKillName(key string) string {
	return strings.TrimPrefix(key, JOB_KILL_DIR)
}

func BuildJobEvent(etype int, jinfo *Job) *JobEvent {
	return &JobEvent{
		EventType: etype,
		JobInfo: jinfo,
	}
}

func BuildJobExecInfo(jsp *JobSchedulePlan) *JobExecInfo {
	ctx, cancel := context.WithCancel(context.TODO())
	return &JobExecInfo{
		Job: jsp.Job,
		PlanTime: jsp.NextTime,
		ScheduleTime: time.Now(),
		Ctx: ctx,
		Cancel: cancel,
	}
}

func BuildJobExecResult(info *JobExecInfo, output string, err error, startTime, endTime time.Time) *JobExecResult {
	return &JobExecResult{
		Info: info,
		Output: output,
		Err: err,
		StartTime: startTime,
		EndTime: endTime,
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

