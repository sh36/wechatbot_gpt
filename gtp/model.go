package gtp

import (
	"fmt"
	"log"
	"strings"
)

// 模型函数类型
type ConversationFunc func(sender string, msg string) (string, error)

// 模型类型
type ConversationModel struct {
	Name string
	Func ConversationFunc
}

// Minimax模型函数
func MinimaxConversation(sender string, msg string) (string, error) {
	// 模型逻辑
	reply, err := Minimax_conversation(sender, msg)
	return reply, err
}

// Xinghuo模型函数
func XinghuoConversation(sender string, msg string) (string, error) {
	// 模型逻辑
	reply, err := Xinghuo_conversation(sender, msg)
	return reply, err
}

// Xinghuo模型函数
func ErnieBotConversation(sender string, msg string) (string, error) {
	// 模型逻辑
	reply, err := ErnieBotApi_conversation(sender, msg)
	return reply, err
}

// 初始化模型
var minimaxModel = ConversationModel{Name: "minimax", Func: MinimaxConversation}
var xinghuoModel = ConversationModel{Name: "星火", Func: XinghuoConversation}
var ErnieBotmodel = ConversationModel{Name: "文心", Func: ErnieBotConversation}

// 当前使用的模型
var currentModel = ErnieBotmodel

var count = 0

func Completions(sender string, msg string) (string, error) {
	currentModel = ErnieBotmodel
	// 判断是否切换模型
	if strings.HasPrefix(msg, "minimax") {
		currentModel = minimaxModel
		msg = strings.TrimSpace(strings.ReplaceAll(msg, "minimax", ""))
		//msg = strings.TrimPrefix(msg, "minimax")
	} else if strings.HasPrefix(msg, "星火") {
		currentModel = xinghuoModel
		msg = strings.TrimSpace(strings.ReplaceAll(msg, "星火", ""))
		//msg = strings.TrimPrefix(msg, "星火")
	} else if strings.HasPrefix(msg, "文心") {
		currentModel = ErnieBotmodel
		msg = strings.TrimSpace(strings.ReplaceAll(msg, "文心", ""))
	}
	// 调用当前模型进行智能交互
	log.Println(msg)
	reply, err := currentModel.Func(sender, msg)
	if err != nil {
		return "", fmt.Errorf("模型调用出错：%s", err.Error())
	}
	count = count + 1
	// 构造回复结果
	log.Printf("当前已调用次数: %d \n", count)

	result := fmt.Sprintf("%s\n\n——\n当前回复来自于%s，以上是模型生成结果，不代表任何人观点。\n可在提问前输入模型名称切换，如：minimax/星火/文心+问题。", reply, currentModel.Name)

	return result, nil
}
