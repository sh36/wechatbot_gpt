package main

import "github.com/869413421/wechatbot/gtp"

/*
func main() {
	bootstrap.Run()
}
*/

func main() {

	gtp.Completions("a", "5+7=?")
}
