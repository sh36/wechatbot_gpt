package gtp

import (
	"fmt"
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

// 初始化模型
var minimaxModel = ConversationModel{Name: "minimax", Func: MinimaxConversation}
var xinghuoModel = ConversationModel{Name: "xinghuo", Func: XinghuoConversation}

// 当前使用的模型
var currentModel = minimaxModel

func Completions(sender string, msg string) (string, error) {
	// 判断是否切换模型
	if strings.HasPrefix(msg, "minimax") {
		currentModel = minimaxModel
		msg = strings.TrimPrefix(msg, "minimax")
	} else if strings.HasPrefix(msg, "xinghuo") {
		currentModel = xinghuoModel
		msg = strings.TrimPrefix(msg, "xinghuo")
	}
	// 调用当前模型进行智能交互
	reply, err := currentModel.Func(sender, msg)
	if err != nil {
		return "", fmt.Errorf("模型调用出错：%s", err.Error())
	}

	// 构造回复结果

	result := reply + fmt.Sprintf("\n当前回复来自于%s，以上是模型生成结果，不代表任何人观点。\n可在提问前输入模型名称切换，如：xinghuo/minimax+问题。", currentModel.Name)

	return result, nil
}
