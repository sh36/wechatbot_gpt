package gtp

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type WsParam struct {
	APPID     string
	APIKey    string
	APISecret string
	GPTURL    string
}

func NewWsParam(APPID, APIKey, APISecret, gptURL string) *WsParam {
	return &WsParam{
		APPID:     APPID,
		APIKey:    APIKey,
		APISecret: APISecret,
		GPTURL:    gptURL,
	}
}

func (w *WsParam) generateSignature() (string, error) {
	now := time.Now().UTC().Format(time.RFC1123)

	signatureOrigin := fmt.Sprintf("host: %s\n", w.getHost())
	signatureOrigin += fmt.Sprintf("date: %s\n", now)
	signatureOrigin += fmt.Sprintf("GET %s HTTP/1.1", w.getPath())

	h := hmac.New(sha256.New, []byte(w.APISecret))
	_, err := h.Write([]byte(signatureOrigin))
	if err != nil {
		return "", err
	}
	signatureSHA := h.Sum(nil)
	signatureSHABase64 := base64.StdEncoding.EncodeToString(signatureSHA)

	authorizationOrigin := fmt.Sprintf(`api_key="%s", algorithm="hmac-sha256", headers="host date request-line", signature="%s"`, w.APIKey, signatureSHABase64)
	authorization := base64.StdEncoding.EncodeToString([]byte(authorizationOrigin))

	return authorization, nil
}

func (w *WsParam) getHost() string {
	u, _ := url.Parse(w.GPTURL)
	return u.Host
}

func (w *WsParam) getPath() string {
	u, _ := url.Parse(w.GPTURL)
	return u.Path
}

func (w *WsParam) createURL() (string, error) {
	authorization, err := w.generateSignature()
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Set("authorization", authorization)
	params.Set("date", time.Now().UTC().Format(time.RFC1123))
	params.Set("host", w.getHost())

	url := w.GPTURL + "?" + params.Encode()
	return url, nil
}

func genParams(appID, question string) ([]byte, error) {
	data := map[string]interface{}{
		"header": map[string]interface{}{
			"app_id": appID,
			"uid":    "1234",
		},
		"parameter": map[string]interface{}{
			"chat": map[string]interface{}{
				"domain":           "general",
				"random_threshold": 0.5,
				"max_tokens":       2048,
				"auditing":         "default",
			},
		},
		"payload": map[string]interface{}{
			"message": map[string]interface{}{
				"text": []map[string]string{
					{"role": "user", "content": question},
				},
			},
		},
	}

	params, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return params, nil
}

func Xinghuo_chat(msg string) (string, error) {
	appID := "89ed2ea4"
	apiKey := "0736d4e16d9cd2bb55d4b1bd38d07229"
	apiSecret := "MDlkZTRlYzNlZTcwZmU5YjU3ZmVmNWNj"
	gptURL := "wss://spark-api.xf-yun.com/v1.1/chat"
	question := msg
	var reply string
	wsParam := NewWsParam(appID, apiKey, apiSecret, gptURL)
	url, err := wsParam.createURL()
	if err != nil {
		log.Println("createURL error:", err)
		return reply, err
	}

	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Println("WebSocket dial error:", err)
		return reply, err
	}
	defer c.Close()

	data, err := genParams(appID, question)
	if err != nil {
		log.Println("genParams error:", err)
		return reply, err
	}

	err = c.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Println("WebSocket write error:", err)
		return reply, err
	}

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			return reply, err
		}

		var response map[string]interface{}
		err = json.Unmarshal(message, &response)
		if err != nil {
			log.Println("JSON unmarshal error:", err)
			return reply, err
		}

		header := response["header"].(map[string]interface{})
		code := int(header["code"].(float64))

		if code != 0 {
			log.Printf("请求错误: %d, %v\n", code, response)
		} else {
			payload := response["payload"].(map[string]interface{})
			choices := payload["choices"].(map[string]interface{})
			status := int(choices["status"].(float64))
			content := choices["text"].([]interface{})[0].(map[string]interface{})["content"].(string)
			reply += content
			//fmt.Print(reply)
			if status == 2 {
				//log.Println("Connection closed")
				break
			}
		}
	}
	return reply, err
}

func Xinghuo_conversation(sender string, msg string) (string, error) {
	// 读取存储的历史记录
	// TODO 根据wx_id获取历史对话

	reply, err := Xinghuo_chat(msg)
	return reply, err
}
