package main

import (
	"github.com/jasonlvhit/gocron"
	"gocn/config"
	"gocn/db"
	"gocn/dingding"
	"gocn/splider"
	"time"
)

func main() {
	if !config.Config.GetBool("splider.all") {
		if config.Config.GetBool("cron.start") {
			runTime := config.Config.GetString("cron.runTime")
			gocron.Every(1).Day().At(runTime).Do(run)
			<-gocron.Start()
			return
		}
	}
	run()
}

func run() {
	go db.Run()
	go dingding.Send()
	defer db.Close()

	splider.Run()

	// 当使用定时任务启动时，使用这里，等待输入写入文件完成
	time.Sleep(time.Minute * 1)
}
