package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	dingtalk "github.com/hugozhu/godingtalk"
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
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// log.Printf("%v", req)
		action := req.URL.Query().Get("action")
		if action == "ding" {
			text := req.URL.Query().Get("text")
			log.Println("text: ", text)
			if text != "" {
				text = "/echo " + text
				exec.Command("/bin/bash", workDirPath+"/cmd.sh", text)
			}
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
	})

	// One can use generate_cert.go in crypto/tls to generate cert.pem and key.pem.
	go func() {
		log.Printf("About to listen on 8888. Go to http://127.0.0.1:8888/")
		err2 := http.ListenAndServe(":8888", nil)
		log.Fatal(err2)
	}()

	log.Printf("About to listen on 8443. Go to https://127.0.0.1:8443/")
	err := http.ListenAndServeTLS(":8443", workDirPath+"/cert.pem", workDirPath+"/key.pem", nil)
	log.Fatal(err)
}
