package main

import (
	"flag"
	"log"
	"os"
	"regexp"

	dingtalk "github.com/hugozhu/godingtalk"
)

var img *bool
var f *string

func init() {
	img = flag.Bool("img", false, "Sending Image")
	f = flag.String("f", "", "File Path")
	flag.Parse()
}

func stripeMarkdown(str string) string {
	str = regexp.MustCompile("[*|#]+").ReplaceAllString(str, "")
	str = regexp.MustCompile("\\s+").ReplaceAllString(str, " ")
	str = regexp.MustCompile("^ ").ReplaceAllString(str, "")
	return str
}

func main() {
	c := dingtalk.NewDingTalkClient(os.Getenv("corpid"), os.Getenv("corpsecret"))
	c.RefreshAccessToken()
	var title string
	var message string
	if *img {
		f, err := os.Open(*f)
		defer f.Close()
		if err != nil {
			log.Fatal(err)
		}
		m, err := c.UploadMedia("image", "screenshot.jpg", f)
		if err != nil {
			log.Fatal(err)
		}
		//log.Printf("%v %v %v", m.MediaID, *f, err)
		markdown := "![Screenshot](" + m.MediaID + ")"
		title = "Screenshot"
		message = markdown
	} else {
		title = stripeMarkdown(os.Args[1])
		message = os.Args[1]
	}
	resp, err := c.SendRobotMarkdownMessage(os.Getenv("token"), title, message)

	if err!=nil {
		panic(err)
	} else {
		log.Println(resp)
	}	
}
