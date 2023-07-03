package handlers

import (
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/869413421/wechatbot/gtp"
	"github.com/eatmoreapple/openwechat"
)

var _ MessageHandlerInterface = (*UserMessageHandler)(nil)

// UserMessageHandler 私聊消息处理
type UserMessageHandler struct {
}

// handle 处理消息
func (g *UserMessageHandler) handle(msg *openwechat.Message) error {
	if msg.IsText() {
		return g.ReplyText(msg)
	}
	return nil
}

// NewUserMessageHandler 创建私聊处理器
func NewUserMessageHandler() MessageHandlerInterface {
	return &UserMessageHandler{}
}

// ReplyText 发送文本消息到群
func (g *UserMessageHandler) ReplyText(msg *openwechat.Message) error {
	// 接收私聊消息
	sender, err := msg.Sender()
	log.Printf("Received User %v Text Msg : %v", sender.NickName, msg.Content)

	// 向GPT发起请求
	requestText := strings.TrimSpace(msg.Content)
	requestText = strings.Trim(msg.Content, "\n")
	reply, err := gtp.Completions(sender.NickName, requestText)
	if err != nil {
		log.Printf("gtp request error: %v \n", err)
		msg.ReplyText("很抱歉，此问题无法回答，请稍后再问。")
		return err
	}
	if reply == "" {
		return nil
	}

	// 回复用户
	reply = strings.TrimSpace(reply)
	reply = strings.Trim(reply, "\n")
	_, err = msg.ReplyText(reply)
	if err != nil {
		log.Printf("response user error: %v \n", err)
	}
	return err
}

// 每隔8小时向指定人发送一条消息
func SendToPerson(bot *openwechat.Bot, userID string, message string) {
	self, err := bot.GetCurrentUser()
	if err != nil {
		log.Println("Error getting current user:", err)
		return
	}

	friends, err := self.Friends()
	if err != nil {
		log.Println("Error getting friends list:", err)
		return
	}

	for {

		result := friends.SearchByRemarkName(1, userID)

		// 在好友列表中查找指定用户
		if len(result) == 0 {
			log.Printf("Could not find user with remark name '%s'\n", userID)
			continue
		}

		// 向指定用户发送一条消息
		if err := result.SendText(message); err != nil {
			log.Printf("Error sending message to '%s': %v\n", userID, err)
		} else {
			log.Printf("Sent message '%s' to %s at %s\n", message, userID, time.Now().Format("2023-04-14 16:13:36"))
		}

		// 生成发送消息的随机时间点
		interval := time.Duration(rand.Intn(5)) * time.Minute
		duration := 8*time.Hour + interval

		// 等待随机时间点
		log.Printf("Waiting for %v before sending next message to %s\n", duration, userID)
		<-time.After(duration)

	}
}
