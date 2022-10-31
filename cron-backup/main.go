package main

import (
	"fmt"
	"time"

	cron "github.com/robfig/cron/v3"
)

// https://segmentfault.com/a/1190000039647260
func main() {
	// c := cron.New(cron.WithSeconds())
	// // 含义查看下文表达式示例
	// c.AddFunc("0/7 * * * * *", func() {
	// 	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	// })
	// c.Start()
	// time.Sleep(300 * time.Second)

	fmt.Println("每天凌晨3点自动备份，不要关！")
	c := cron.New(cron.WithSeconds())
	spec := "0/7 * * * * *"
	c.AddFunc(spec, func() {
		inixcopy()
	})

	c.AddFunc("0/5 * * * * *", func() {
		fmt.Println(".......")
	})

	c.Start()
	//time.Sleep(300 * time.Second)
	select {} //阻塞主线程停止
}

func inixcopy() {
	fmt.Println("被执行了")
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}
