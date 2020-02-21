package main

import (
	"os"
    "regexp"

	dingtalk "github.com/hugozhu/godingtalk"
)

func stripeMarkdown(str string) string {
   	str = regexp.MustCompile("[*|#]+").ReplaceAllString(str, "")
	str = regexp.MustCompile("\\s+").ReplaceAllString(str, " ")
	str = regexp.MustCompile("^ ").ReplaceAllString(str, "")
	return str
}

func main() {
	c := dingtalk.NewDingTalkClient(os.Getenv("corpid"), os.Getenv("corpsecret"))
	c.RefreshAccessToken()
	c.SendRobotMarkdownMessage(os.Getenv("token"), stripeMarkdown(os.Args[1]), os.Args[1])
}
