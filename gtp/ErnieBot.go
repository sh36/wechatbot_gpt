package gtp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const ErnieBotURL = "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/eb-instant"

type ErnieBotResponseBody struct {
	ID               string `json:"id"`
	Object           string `json:"object"`
	Created          int    `json:"created"`
	SentenceID       int    `json:"sentence_id"`
	IsEnd            bool   `json:"is_end"`
	Result           string `json:"result"`
	NeedClearHistory bool   `json:"need_clear_history"`
}

type ErnieBotMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ErnieBotRequestBody struct {
	Messages []ErnieBotMessage `json:"messages"`
	Stream   bool              `json:"stream,omitempty"`
	UserID   string            `json:"user_id,omitempty"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type ErnieBotChatHistory struct {
	History       []ErnieBotMessage `json:"history"`
	Conversations int               `json:"clear"`
}

var chatHistory = ErnieBotChatHistory{
	History:       []ErnieBotMessage{},
	Conversations: 0,
}

func GetAccessToken(clientID, clientSecret string) (string, error) {
	url := "https://aip.baidubce.com/oauth/2.0/token"
	payload := strings.NewReader(fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", clientID, clientSecret))
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	response := &AccessTokenResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return "", err
	}
	log.Printf("response.AccessToken: %s\n", response.AccessToken)

	return response.AccessToken, nil
}

// ErnieBot文本模型回复
func ErnieBot_chat(msg string) (string, error) {
	accessToken, err := GetAccessToken("aX9xYA7eQi9nvMF2cRwyDG0q", "NrWvvEPBIeqLRwridSr3RUqtd5CZhUcA")
	if err != nil {
		return "", err
	}

	requestBody := ErnieBotRequestBody{
		Messages: []ErnieBotMessage{
			{
				Role:    "user",
				Content: msg,
			},
		},
		Stream: false, // 设置stream的值，如果不需要使用流式接口则为false
		UserID: "",    // 设置user_id的值，如果不需要指定用户ID则为空字符串
	}

	/*
		// 构建请求体
		requestBody := ErnieBotRequestBody{
			Messages: chatHistory.History,
			Stream:   false, // 设置stream的值，如果不需要使用流式接口则为false
			UserID:   "",    // 设置user_id的值，如果不需要指定用户ID则为空字符串
		}

		// 添加用户输入到请求体中
		requestBody.Messages = append(requestBody.Messages, ErnieBotMessage{
			Role:    "user",
			Content: msg,
		})
	*/
	requestData, err := json.Marshal(requestBody)

	if err != nil {
		return "", err
	}
	fmt.Printf("request gtp json string: %v\n", string(requestData))
	newErnieBotURL := ErnieBotURL + "?access_token=" + accessToken
	req, err := http.NewRequest("POST", newErnieBotURL, bytes.NewBuffer(requestData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

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
	ernieBotResponseBody := &ErnieBotResponseBody{}
	err = json.Unmarshal(body, ernieBotResponseBody)
	if err != nil {
		return "", err
	}

	reply := ernieBotResponseBody.Result
	fmt.Printf("ernieBotResponseBody json string: %v\n", ernieBotResponseBody)
	if reply == "" {
		reply = "很抱歉，此问题无法回答，请稍后再问。"
	}
	log.Printf("gpt response text: %s \n", reply)

	return reply, nil
}

func ErnieBot_conversation(sender string, msg string) (string, error) {
	// 读取存储的历史记录
	// TODO 根据wx_id获取历史对话
	// 将本次对话的问题和答案存入历史记录
	// if chatHistory.Clear {
	// 	chatHistory.History = []ErnieBotMessage{}
	// 	chatHistory.Clear = false
	// }
	/*
		chatHistory.Conversations += 1
		if chatHistory.Conversations >= 5 {
			chatHistoryClear()
		}
	*/

	reply, err := ErnieBot_chat(msg)
	if err != nil {
		return "", err
	}
	/*
		chatHistory.History = append(chatHistory.History, ErnieBotMessage{
			Role:    "user",
			Content: msg,
		})
		chatHistory.History = append(chatHistory.History, ErnieBotMessage{
			Role:    "assistant",
			Content: reply,
		})
	*/
	return reply, err
}

func chatHistoryClear() {
	chatHistory.History = []ErnieBotMessage{}
	chatHistory.Conversations = 0
}
