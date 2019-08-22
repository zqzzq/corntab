package master

import (
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"go.mongodb.org/mongo-driver/mongo/options"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/bson"
	"corntab/src/common"
)

type LogSink struct {
	Client *mongo.Client
	Collection *mongo.Collection
}

var (
	S_logSink *LogSink
)

func (ls *LogSink)GetJobLogs(name string) (logList []common.JobLog, err error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(S_config.MongodbTimeout)*time.Second)
	cursor, err := S_logSink.Collection.Find(ctx, bson.D{{"jobName",name}})
	if err != nil{
		fmt.Println("获取任务日志失败：", err)
		return
	}
	defer cursor.Close(ctx)
	logList = []common.JobLog{}
	for cursor.Next(ctx) {
		var result common.JobLog
		err = cursor.Decode(&result)
		if err != nil {
			fmt.Println("解码任务日志失败：", err)
			continue
		}
		logList = append(logList, result)
	}
	return
}

func InitLogSink() (err error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(S_config.MongodbTimeout)*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(S_config.MongodbUrl))
	if err != nil{
		fmt.Println("创建MongoDB client失败：", err)
		return
	}
	if err = client.Ping(ctx, readpref.Primary()); err != nil{
		fmt.Println("连接MongoDB失败：", err)
		return
	}
	S_logSink = &LogSink{
		Client: client,
		Collection: client.Database("cron").Collection("log"),
	}
	return
}