package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type OutgoingMessage struct {
	MessageType string `json:"msgtype"`
	Text        struct {
		Content string `json:"content,omitempty"`
	} `json:"text,omitempty"`
	MessageID         string `json:"msgId"`
	CreatedAt         int64  `json:"createAt"`
	ConversationID    string `json:"conversationId"`
	ConversationType  string `json:"conversationType"`
	ConversationTitle string `json:"conversationTitle"`
	SenderID          string `json:"senderId"`
	SenderNick        string `json:"senderNick"`
	SenderCorpID      string `json:"senderCorpId"`
	SenderStaffID     string `json:"senderStaffId"`
	ChatbotUserID     string `json:"chatbotUserId"`
	AtUsers           []struct {
		DingTalkID string `json:"dingtalkId,omitempty"`
		StaffID    string `json:"staffId,omitempty"`
	} `json:"atUsers,omitempty"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("%v", req)
		content, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(content))
		var obj OutgoingMessage
		json.Unmarshal(content, &obj)
		log.Printf("%v", obj)
		io.WriteString(w, "Hello, TLS!\n")
	})

	// One can use generate_cert.go in crypto/tls to generate cert.pem and key.pem.
	go func() {
		log.Printf("About to listen on 8888. Go to http://127.0.0.1:8888/")
		err2 := http.ListenAndServe(":8888", nil)
		log.Fatal(err2)
	}()

	log.Printf("About to listen on 8443. Go to https://127.0.0.1:8443/")
	err := http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil)
	log.Fatal(err)
}
