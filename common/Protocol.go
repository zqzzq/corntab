package common

import (
	"encoding/json"
	"strings"
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