package master

import (
	"net/http"
	"net"
	"time"
	"strconv"
	"fmt"
	"corntab/common"
	"encoding/json"
)

//任务的HTTP服务
type ApiServer struct {
	httpServer *http.Server
}
//任务保存方法
//post job={"name":"job1","command":"echo hello","cronExpr":"* * * * *"}
func handleJobSave(resp http.ResponseWriter, req *http.Request)  {
	if err := req.ParseForm();err != nil{
		fmt.Println(err)
		return
	}
	postJob := req.PostForm.Get("job")
	var job common.Job
	if err :=  json.Unmarshal([]byte(postJob), &job);err !=nil{
		fmt.Println(err)
		return
	}
	//保存到etcd
	oldJob, err := S_jobMgr.SaveJob(&job)
	if err != nil{
		fmt.Println(err)
		return
	}
	by, err := common.BuildResponse(0, "success", *oldJob)
	if err != nil{
		fmt.Println(err)
		return
	}
	resp.Write(by)
}

//初始化服务
func InitApiServer() (err error) {
	//配置路由
	mux := http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)

	lis, err := net.Listen("tcp", ":" + strconv.Itoa(S_config.ApiPort))
	if err != nil{
		return
	}

	serve := http.Server{
		ReadTimeout: time.Duration(S_config.ApiReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(S_config.ApiWriteTimeout) * time.Millisecond,
		Handler: mux,
	}

	if err = serve.Serve(lis);err != nil{
		return
	}
	return
}