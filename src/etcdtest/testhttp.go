package main

import (
	"time"
	"bytes"
	"io"
	"net/http"
	"fmt"
	"errors"
)

func Get(url string) (response string) {
	client := http.Client{Timeout: 5 * time.Second}
	resp, error := client.Get(url)
	if error != nil {
		fmt.Println(error)
		return
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			return
		}
	}

	response = result.String()
	return
}

func main()  {


	err := errors.New("aaaaaaaa")
	fmt.Println(err)

	//url := "http://liquibase-server.cicd:80/"
	//resp := Get(url)
	//fmt.Println(resp)
	////time.Sleep(180)
	//lis, err := net.Listen("tcp", ":9901")
	//if err != nil{
	//	fmt.Println(err)
	//}
	//server := http.Server{}
	//server.Serve(lis)
}
