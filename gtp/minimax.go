package gtp

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const MinimaxURL = "https://api.minimax.chat/v1/text/chatcompletion?GroupId=1683775659483837"

type MiniMaxResponseBody struct {
	Reply string `json:"reply"`
}

type ChoiceItem struct {
}

type MiniMaxRequestBody struct {
	Model            string        `json:"model"`
	TokensToGenerate int           `json:"tokens_to_generate"`
	Messages         []ChatMessage `json:"messages"`
}

type ChatMessage struct {
	SenderType string `json:"sender_type"`
	Text       string `json:"text"`
}

// Minimax文本模型回复
func Minimax_chat(msg string, history_stack *History_stack) (string, error) {
	history_stack.check_rounds()

	requestBody := MiniMaxRequestBody{
		Model:            "abab5-chat",
		TokensToGenerate: 2048,
		Messages: []ChatMessage{
			{
				SenderType: "USER",
				Text:       msg,
			},
		},
	}

	requestData, err := json.Marshal(requestBody)

	if err != nil {
		return "", err
	}
	log.Printf("request gtp json string : %v", string(requestData))
	req, err := http.NewRequest("POST", MinimaxURL, bytes.NewBuffer(requestData))
	if err != nil {
		return "", err
	}

	//apiKey := config.LoadConfig().ApiKey
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJTdWJqZWN0SUQiOiIxNjgzNzc1NjU5MDg2MzY0IiwiUGhvbmUiOiIxNTcqKioqMzAxMCIsIkdyb3VwSUQiOiIiLCJQYWdlTmFtZSI6IiIsIk1haWwiOiJsaW5ib2hvbmdAY2liLmNvbS5jbiIsImlzcyI6Im1pbmltYXgifQ.LNBQJzKamC6yx7XkPCvURixG_b8X9VoHb1qoBAsoKyFhaoO-xpKl0cpisyxsPcOWBzcG0Ltq3UdgG1cVlUGhFfAH0NHn4h5ZPDW8579aLPLuKomHhVMaxF54JHkPJ0Hzx7_py5ipNhh0e5PDdweit-FEZhhtSkRZ6H4OmnjMTfTi4PCEytVILe6blOUbCQWqZr1bKiblGA3rg10DsZLQsilesiNmeDn3zsXaoU4wPgif8AcNJES8Hh_cknf0rC8J4A4kXZbKKzH0XGJ9am0TyY4jnNwF5rAFbrCJJAceBkCfIukHXdrOB1_Khy5GuzJLhgroLRqsZE487mU3oDmFvA")
	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	gptResponseBody := &MiniMaxResponseBody{}
	log.Println(string(body))
	err = json.Unmarshal(body, gptResponseBody)
	if err != nil {
		return "", err
	}
	var reply string
	if len(gptResponseBody.Reply) > 0 {
		reply = gptResponseBody.Reply
	} else {
		reply = "很抱歉，此问题无法回答，请稍后再问。"
	}
	log.Printf("gpt response text: %s \n", reply)

	*history_stack.History = append(*history_stack.History,
		ChatMessage{
			SenderType: "BOT",
			Text:       reply,
		})

	return reply, nil
}

func Minimax_conversation(sender string, msg string) (string, error) {
	// 读取存储的历史记录
	// TODO 根据wx_id获取历史对话

	history, err := GetHistoryStack(sender)
	if err != nil {
		return "", err
	}
	reply, err := Minimax_chat(msg, history)
	return reply, err
}
