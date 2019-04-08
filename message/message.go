package message

import (
	"errors"
	"fmt"
	"strings"
)

// TextUrl 每日新闻内容
type TextUrl struct {
	Text string
	Url  string
}

type Message struct {
	TextUrls   []TextUrl
	DailyTitle string
	Author     string // 编辑
	PostUrl    string // 原文地址
}

var messageChain chan Message

func init() {
	messageChain = make(chan Message, 10000)
}

func Push(m Message) {
	messageChain <- m
}

func Pop() (Message, error) {
	select {
	case m := <-messageChain:
		return m, nil
	default:
		return Message{}, errors.New("nil")
	}
}


func BuildMessage(msg Message) string {
	str := fmt.Sprintf("## %s", msg.DailyTitle)
	str += fmt.Sprintln()
	for _, v := range msg.TextUrls {
		if len(v.Url) == 0 || len(v.Text) == 0 {
			continue
		}

		if strings.Contains(v.Text, "GoCN归档") || strings.Contains(v.Text, "订阅新闻") {
			continue
		}

		textValue := strings.Replace(v.Text, "\n", "", -1)
		realText := strings.Replace(textValue, " ", "", -1)

		str += fmt.Sprintf("- [%s](%s)", realText, v.Url)
		str += fmt.Sprintln()
	}

	index := strings.Index(msg.Author, "订阅新闻")
	author := msg.Author
	if index > 0 {
		author = msg.Author[:index]
	}

	str += fmt.Sprintln()
	str += fmt.Sprintf("编辑：%s", author)
	str += fmt.Sprintln()
	str += fmt.Sprintln()
	str += fmt.Sprintf("原文地址: %s", msg.PostUrl)
	return str
}
