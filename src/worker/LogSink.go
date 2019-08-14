package worker

import (
	"go.mongodb.org/mongo-driver/mongo"
	"corntab/src/common"
	"time"
	"go.mongodb.org/mongo-driver/mongo/options"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type LogSink struct {
	Client *mongo.Client
	Collection *mongo.Collection
	LogChan chan *common.JobLog
}

var (
	S_logSink *LogSink
)

func (logSink *LogSink) saveLog(batch *common.JobLogBatch)  {
	if _, err :=  logSink.Collection.InsertMany(context.TODO(), batch.JobLogs);err != nil{
		fmt.Println("写日志到MongoDB错误：",err)
	}
}

func (logSink *LogSink) logWriteLoop()  {
	var log *common.JobLog
	var jobLogBatch *common.JobLogBatch
	logCommitTimer := time.NewTimer(time.Duration(S_config.LogCommitTime) * time.Second)
	for {
		select {
		case log =<- logSink.LogChan:
			if jobLogBatch == nil{
				jobLogBatch = &common.JobLogBatch{}
			}
			jobLogBatch.JobLogs = append(jobLogBatch.JobLogs, log)
			if len(jobLogBatch.JobLogs) >= S_config.JobLogBatchSize{
				S_logSink.saveLog(jobLogBatch)
				jobLogBatch = nil
			}
		case <- logCommitTimer.C:
			if jobLogBatch != nil && len(jobLogBatch.JobLogs) != 0 {
				S_logSink.saveLog(jobLogBatch)
				jobLogBatch = nil
			}
			logCommitTimer = time.NewTimer(time.Duration(S_config.LogCommitTime) * time.Second)

		}
	}
}

func (logSink *LogSink)AppendJobLog(log *common.JobLog)  {
	logSink.LogChan <- log
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
		LogChan: make(chan *common.JobLog, 1000),
	}
	go S_logSink.logWriteLoop()
	return
}