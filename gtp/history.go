package gtp

import (
	"errors"
	"fmt"
	"time"

	"github.com/869413421/wechatbot/config"
)

// TODO: 改造history为map，以wx_id为键
type History_stack struct {
	Message_sender string
	History        *[]ChatMessage `json:"messages"`
	Max_boxes      int
}

var Max_boxes = config.LoadConfig().Max_boxes
var User_count = config.LoadConfig().User_count

func New_History_stack(sender string, history *[]ChatMessage, max_boxes int) *History_stack {
	history_stack := History_stack{
		Message_sender: sender,
		History:        history,
		Max_boxes:      max_boxes,
	}
	return &history_stack
}

func (h *History_stack) clear() {
	*h.History = (*h.History)[:0]
}

func (h *History_stack) check_rounds() {
	len := h.count()
	if len > h.Max_boxes {
		// 删除历史
		h.clear()
		//return NewError("轮次过多，已清空上下文。请重新提问。")
	}
	return
}

func (h *History_stack) count() int {
	return len(*h.History)
}

var History_stack_slice []*History_stack

// GetHistoryStack 函数根据 Message_sender 查找历史栈，如果不存在则创建一个新的，并返回 History_stack 对象
func GetHistoryStack(sender string) (*History_stack, error) {
	for _, stack := range History_stack_slice {
		if stack.Message_sender == sender {
			return stack, nil
		}
	}

	//判断当前用户数量，控制并发数
	if len(History_stack_slice) >= User_count {
		return nil, errors.New("exceeded maximum number of users")
	}
	// 如果历史栈不存在，则创建一个新的
	historyStack := New_History_stack(sender, &[]ChatMessage{}, Max_boxes)
	History_stack_slice = append(History_stack_slice, historyStack)
	return historyStack, nil
}

// 定时清空历史记录
func ClearHistoryStackSlice() {
	/*
	   for {
	       time.Sleep(48 * time.Hour) // 等待 48 小时
	       History_stack_slice = []*History_stack{} // 清空历史记录切片
	   }
	*/
	for {
		// 计算距离下一个早上三点的时间
		now := time.Now()
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 3, 0, 0, 0, next.Location())
		duration := next.Sub(now)

		// 等待到下一个早上三点
		time.Sleep(duration)

		// 清空历史记录切片
		History_stack_slice = []*History_stack{}
		fmt.Printf("Cleared history stack slice at %s\n", time.Now().Format("2006-01-02 15:04:05"))
	}
}

type TooMuchRound struct {
	msg string
}

func (e TooMuchRound) Error() string {
	return fmt.Sprintf("msg:%v", e.msg)
}

func NewError(msg string) error {
	return TooMuchRound{
		msg: msg,
	}
}

func (e TooMuchRound) GetMessage() string {
	return e.msg
}
