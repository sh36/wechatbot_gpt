package gtp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func GetAccessToken(clientID, clientSecret string) (string, error) {
	url := fmt.Sprintf("https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/eb-instant")
	payload := strings.NewReader("")
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

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
	fmt.Printf("response.AccessToken: %s\n", response.AccessToken)

	return response.AccessToken, nil

}

// ErnieBot文本模型回复
func ErnieBot_chat(msg string, history []ErnieBotMessage) (string, error) {
	accessToken, err := GetAccessToken("aX9xYA7eQi9nvMF2cRwyDG0q", "NrWvvEPBIeqLRwridSr3RUqtd5CZhUcA")
	if err != nil {
		return "", err
	}
	requestBody := ErnieBotRequestBody{
		Messages: append(history, ErnieBotMessage{
			Role:    "user",
			Content: msg,
		}),
		Stream: false, // 设置stream的值，如果不需要使用流式接口则为false
		UserID: "",    // 设置user_id的值，如果不需要指定用户ID则为空字符串
	}

	requestData, err := json.Marshal(requestBody)

	fmt.Printf("request gpt json string: %v", string(requestData))

	if err != nil {
		return "", err
	}
	//log.Printf("request gtp json string: %v", string(requestData))
	fmt.Printf("request gtp json string: %v\n", string(requestData))
	req, err := http.NewRequest("POST", ErnieBotURL, bytes.NewBuffer(requestData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("access_token", accessToken)
	req.Header.Set("Authorization", "Bearer "+accessToken)

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

	ernieBotResponseBody := &ErnieBotResponseBody{}
	err = json.Unmarshal(body, ernieBotResponseBody)
	if err != nil {
		return "", err
	}

	reply := ernieBotResponseBody.Result
	if reply == "" {
		reply = "很抱歉，此问题无法回答，请稍后再问。"
	}
	fmt.Printf("gpt response text: %s \n", reply)

	return reply, nil
}

func ErnieBot_conversation(sender string, msg string) (string, error) {
	// 读取存储的历史记录
	// TODO 根据wx_id获取历史对话

	history := make([]ErnieBotMessage, 0)
	reply, err := ErnieBot_chat(msg, history)
	return reply, err
}
