package main

import (
	"google.golang.org/grpc"
	"test/rpc/pb"
	"time"
	"log"
	"context"
)

const (
	address     = "localhost:50053"
)

func main() {
	i := 0
	for {

		if i < 5 {
			go func(add string, t time.Duration) {
				conn, err := grpc.Dial(add, grpc.WithInsecure())
				if err != nil {
					log.Fatalf("did not connect: %v", err)
				}
				defer conn.Close()
				c := pb.NewExecuteServiceClient(conn)

				ctx, cancel := context.WithTimeout(context.Background(), t * time.Second)
				defer cancel()
				resp, err := c.Execute(ctx, &pb.ExecuteRequest{})
				if err != nil {
					log.Fatalf("执行失败: %v", err)

				}
				log.Println("执行结果是：", resp.Result)
			}(address, 4)
		}else {
			log.Println("执行完毕 退出")
		}
		time.Sleep( 2 * time.Second )
		i ++
	}
}