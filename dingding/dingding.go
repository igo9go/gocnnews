package dingding

import (
	"encoding/json"
	"fmt"
	"gocn/config"
	"gocn/db"
	"gocn/message"
	"time"
)

func Send() {
	ding := Ding{AccessToken: config.Config.GetString("dingding.token")}
	for {
		msg, err := message.Pop()
		if err != nil {
			continue
		}
		if !db.HasSave(msg.DailyTitle) {
			data, err := json.Marshal(msg)
			if err != nil {
				fmt.Printf("json.marshal failed,err:", err)
			}
			db.Push(string(data))
			content := message.BuildMessage(msg)
			db.PushMarkdown(content)
			markdown := Markdown{Title: "GoCN每日新闻", Content: content}
			ding.Send(markdown)
			time.Sleep(time.Second * 1)
		}
	}
}
