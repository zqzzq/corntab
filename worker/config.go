package worker

import (
	"io/ioutil"
	"encoding/json"
	"fmt"
)
//单例
var (
	S_config *Config
)

//程序配置
type Config struct {
	EtcdDialTimeout int `json:"etcdDialTimeout"`
	EtcdEndpoints string `json:"etcdEndpoints"`
}

//加载配置
func InitConfig(filePath string) (err error) {
	content, err := ioutil.ReadFile(filePath)
	if err!=nil {
		return
	}
	conf := Config{}
	if err = json.Unmarshal(content, &conf); err != nil{
		return
	}
	S_config = &conf
	fmt.Println(S_config)
	return
}