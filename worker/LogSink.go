package worker

import (
	"go.mongodb.org/mongo-driver/mongo"
	"corntab/common"
	"time"
	"go.mongodb.org/mongo-driver/mongo/options"
	"context"
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
	logSink.Collection.InsertMany(context.TODO(), batch.JobLogs)
}

func (logSink *LogSink) logWriteLoop()  {
	var log *common.JobLog
	var jobLogBatch *common.JobLogBatch
	for {
		select {
		case log =<- logSink.LogChan:
			if jobLogBatch == nil{
				jobLogBatch = &common.JobLogBatch{}
			}
			jobLogBatch.JobLogs = append(jobLogBatch.JobLogs, log)
			if len(jobLogBatch.JobLogs) >= S_config.JobLogBatchSize{
				S_logSink.saveLog(jobLogBatch)
			}


		}
	}
}

func InitLogSink() (err error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(S_config.MongodbTimeout)*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(S_config.MongodbUrl))
	if err != nil{
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