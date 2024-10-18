package robot

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type FeishuRobot struct {
	Url   string `yaml:"url"`
	Times int    `yaml:"times"`
}

func newMessage(text string) map[string]interface{} {
	message := make(map[string]interface{})
	message["msg_type"] = "text"
	content := map[string]string{"text": text}
	message["content"] = content
	return message
}

func (f *FeishuRobot) Send2robot(text string) {
	messageJSON, err := json.Marshal(newMessage(text))
	if err != nil {
		log.Printf("error %s\n", err)
	}
	req, err := http.NewRequest("POST", f.Url, bytes.NewBuffer(messageJSON))
	if err != nil {
		log.Printf("error %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error %s", err)
	}
	defer resp.Body.Close()
}
