package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/869413421/wechatbot/gtp"
)

/*
func main() {
	bootstrap.Run()
}s
*/

func main() {
	// 循环测试
	for {
		// 询问用户输入
		fmt.Print("请输入测试内容（输入 q 退出）: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// 检查是否退出
		if input == "q" {
			fmt.Println("退出测试")
			break
		}

		// 调用 Completions 函数进行测试
		reply, err := gtp.Completions("sender", input)
		if err != nil {
			fmt.Printf("调用出错：%s\n", err.Error())
			continue
		}

		// 打印回复结果
		fmt.Printf("%s\n", reply)
	}
}
