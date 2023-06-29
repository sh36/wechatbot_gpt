package gtp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/869413421/wechatbot/config"
)

const SenseTimeURL = "https://lm_experience.sensetime.com/v1/nlp/chat/completions"
var SenseTime_api_secret_key = config.LoadConfig().SenseTime_api_secret_key

type SenseTimeRequestBody struct {
	// Model			  string	`json:"model"`
	Messages []SenseTimeMessage `json:"messages"`
	// Temperature       float64 	`json:"temperature"`
	// TopP              float64 	`json:"top_p"`
	// MaxNewTokens      int     	`json:"max_new_tokens"`
	// RepetitionPenalty float64 	`json:"repetition_penalty"`
	// Stream            bool    	`json:"stream"`
	// User              string  	`json:"user"`
}

type SenseTimeResponseBody struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		ID      string `json:"id"`
		Choices []struct {
			Message      string `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
		// 0：成功，1：内部错误，3：超时，4：内容违规
		Status int `json:"status"`
	} `json:"data"`
}

// 在数组中代表history，最后一个为最新的查询消息
type SenseTimeMessage struct{
	Role    string `json:"role"`
	Content string `json:"content"`
}

// SenseTime文本模型回复
func SenseTime_chat(msg string) (string, error) {
	requestBody := SenseTimeRequestBody{
		Messages: []SenseTimeMessage{
			{
				Role:    "user",
				Content: msg,
			},
		},
	}

	requestData, err := json.Marshal(requestBody)

	if err != nil {
		return "", err
	}
	fmt.Printf("request sensetime_lm json string: %v\n", string(requestData))
	req, err := http.NewRequest("POST", SenseTimeURL, bytes.NewBuffer(requestData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", SenseTime_api_secret_key)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	fmt.Printf("response.Body json string: %v\n", response.Body)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	SenseTimeResponseBody := &SenseTimeResponseBody{}
	err = json.Unmarshal(body, SenseTimeResponseBody)
	if err != nil {
		return "", err
	}
	// 目前未考虑流式和长度
	reply := SenseTimeResponseBody.Data.Choices[0].Message
	fmt.Printf("SenseTimeResponseBody json string: %v\n", SenseTimeResponseBody)
	if reply == "" {
		reply = "很抱歉，此问题无法回答，请稍后再问。"
	}
	log.Printf("gpt response text: %s \n", reply)

	return reply, nil
}

func SenseTime_conversation(sender string, msg string) (string, error) {
	reply, err := SenseTime_chat(msg)
	return reply, err
}

