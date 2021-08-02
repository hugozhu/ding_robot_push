package main

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	proto "github.com/hugozhu/ding_robot_push/chatbot/chat" // 根据proto文件自动生成的代码
	"google.golang.org/grpc"
)

var workDirPath string

func init() {
	var err error
	workDirPath, err = os.Getwd()
	if err != nil {
		log.Println(err)
	}
}

func main() {
	url := os.Getenv("grpc_server")
	if url == "" {
		url = "localhost:3000"
	}
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect: [%v]\n", err)
		return
	}
	defer conn.Close()

	// 声明客户端
	client := proto.NewChatClient(conn)

	// 声明 context
	ctx := context.Background()

	// 创建双向数据流
	stream, err := client.BidStream(ctx)
	if err != nil {
		log.Printf("Streaming error: [%v]\n", err)
	}

	go func() {
		for {
			if err := stream.Send(&proto.Request{Input: "/ping"}); err != nil {
				log.Println(err)
			}
			time.Sleep(time.Duration(5) * time.Second)
		}
	}()

	for {
		message, err := stream.Recv()
		if err != nil {
			log.Println("Error:", err)
			break
		}
		log.Printf("[Received]: %s", message.Output)
		text := message.Output
		if strings.HasPrefix(text, "/pong") {

		} else {
			cmd := exec.Command("/bin/bash", workDirPath+"/cmd.sh", text)
			var stderr, stdout bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err = cmd.Run()
			if err != nil {
				log.Println(err.Error(), stderr.String())
			} else {
				log.Println(stdout.String())
			}
		}
	}
}
