package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"

	dingtalk "github.com/hugozhu/godingtalk"

	proto "github.com/hugozhu/ding_robot_push/chatbot/chat" // 自动生成的 proto代码
	"google.golang.org/grpc"
)

var workDirPath string

// Streamer 服务端
type Streamer struct {
	proto.UnimplementedChatServer
}

// 所有连接的客户端
var clients map[string]proto.Chat_BidStreamServer

// 广播消息给所有客户端
func broadCast(message string) error {
	for _, client := range clients {
		if err := client.Send(&proto.Response{Output: message}); err != nil {
			return err
		}
	}
	return nil
}

func (s *Streamer) BidStream(stream proto.Chat_BidStreamServer) error {
	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			sid := fmt.Sprintf("%v", stream)
			delete(clients, sid)
			// log.Println("收到客户端通过context发出的终止信号")
			return ctx.Err()
		default:
			// 接收从客户端发来的消息
			message, err := stream.Recv()
			sid := fmt.Sprintf("%v", &stream)
			if err == io.EOF {
				// log.Println("客户端发送的数据流结束")
				return nil
			}
			if err != nil {
				// log.Println("接收数据出错:", err)
				return err
			}
			clients[sid] = stream
			// 如果接收正常，则根据接收到的 字符串 执行相应的指令
			switch message.Input {
			case "/ping":
				if err := stream.Send(&proto.Response{Output: "/pong"}); err != nil {
					return err
				}
			default:
				log.Printf("[收到消息]: %s", message.Input)
			}
		}
	}
}

func init() {
	var err error
	workDirPath, err = os.Getwd()
	if err != nil {
		log.Println(err)
	}
}

func handler(w http.ResponseWriter, req *http.Request) {
	action := req.URL.Query().Get("action")
	if action == "ding" {
		text := req.URL.Query().Get("text")
		log.Println("text: ", text)
		if text != "" {
			exec.Command("/bin/bash", workDirPath+"/push.sh", text).Run()
		}
	} else if action == "send" {
		text := req.URL.Query().Get("text")
		err := broadCast(text)
		if err != nil {
			io.WriteString(w, "{ \"errcode\": 1, \"errmsg\": \", err.Error() ,\"}")
		} else {
			io.WriteString(w, "{ \"errcode\": 0, \"errmsg\": \"ok\"}")
		}
	} else {
		log.Printf("%v", req)
		content, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(content))
		var obj dingtalk.RobotOutgoingMessage
		json.Unmarshal(content, &obj)
		text := obj.Text.Content
		log.Printf("%s", text)
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
		io.WriteString(w, "{ \"errcode\": 0, \"errmsg\": \"ok\"}")
	}
}

func main() {
	http.HandleFunc("/", handler)

	// One can use generate_cert.go in crypto/tls to generate cert.pem and key.pem.
	go func() {
		log.Printf("Web Server listens on 8888")
		err2 := http.ListenAndServe(":8888", nil)
		log.Fatal(err2)
	}()

	// log.Printf("About to listen on 8443. Go to https://127.0.0.1:8443/")
	// err := http.ListenAndServeTLS(":8443", workDirPath+"/cert.pem", workDirPath+"/key.pem", nil)
	// log.Fatal(err)

	// Start grpc chat server
	server := grpc.NewServer()
	clients = make(map[string]proto.Chat_BidStreamServer)

	// 注册 ChatServer
	proto.RegisterChatServer(server, &Streamer{})
	address, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}
	log.Printf("Chat Server listens on 3000")
	if err := server.Serve(address); err != nil {
		panic(err)
	}
}
