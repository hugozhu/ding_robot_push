module github.com/hugozhu/ding_robot_push/chatbot

go 1.16

replace github.com/hugozhu/ding_robot_push => ../

require (
	github.com/golang/protobuf v1.5.2
	github.com/hugozhu/godingtalk v1.0.5
	google.golang.org/grpc v1.53.0
	google.golang.org/protobuf v1.28.1
)
