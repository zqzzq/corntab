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
	fmt.Println("开始处理job保存")
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
	fmt.Println("保存完成")
	by, err := common.BuildResponse(0, "success", oldJob)
	if err != nil{
		fmt.Println(err)
		return
	}

	resp.Write(by)
}
//任务删除方法
//GET /job/delete?name=job1
func handleJobDelete(resp http.ResponseWriter, req *http.Request)  {
	urlParam := req.URL.Query()
	jobName := urlParam.Get("name")
	fmt.Println("删除job ：", jobName)
	oldJob, err := S_jobMgr.DeleteJob(jobName)
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println("删除完成")
	by, err := common.BuildResponse(0, "success", oldJob)
	if err != nil{
		fmt.Println(err)
		return
	}

	resp.Write(by)

}
//获取所有任务列表
//GET /jobs/list
func handleJobList(resp http.ResponseWriter, req *http.Request)  {
	jobList, err := S_jobMgr.ListJob()
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println("获取任务列表完成")
	by, err := common.BuildResponse(0, "success", jobList)
	if err != nil{
		fmt.Println(err)
		return
	}

	resp.Write(by)
}
//杀死任务
//GET /jobs/kill?name=job1
func handleJobKill(resp http.ResponseWriter, req *http.Request)  {
	urlParam := req.URL.Query()
	jobName := urlParam.Get("name")
	fmt.Println("杀死job ：", jobName)
	oldJob, err := S_jobMgr.KillJob(jobName)
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println("杀死完成")
	by, err := common.BuildResponse(0, "success", oldJob)
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
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list", handleJobList)
	mux.HandleFunc("/job/kill", handleJobKill)

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