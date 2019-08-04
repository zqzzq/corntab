package common

import "encoding/json"

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

func BuildResponse(respNum int, msg string, data interface{}) (resp []byte, err error) {
	r := Response{
		RespNum: respNum,
		Msg: msg,
		Data: data,
	}
	resp, err = json.Marshal(r)
	return
}