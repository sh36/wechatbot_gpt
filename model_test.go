package main

import (
	"fmt"
	"testing"

	"github.com/869413421/wechatbot/gtp"
)

func TestModel(t *testing.T) {
	// 调用 Completions 函数进行测试
	reply, err := gtp.Completions("sender", "你好，你是谁？")
	if err != nil {
		fmt.Printf("调用出错：%s\n", err.Error())
	}

	// 打印回复结果
	fmt.Printf("%s\n", reply)
}
