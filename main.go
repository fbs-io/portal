/*
 * @Author: reel
 * @Date: 2023-06-24 08:50:38
 * @LastEditors: reel
 * @LastEditTime: 2023-09-06 07:29:34
 * @Description: 入口函数
 */
package main

import (
	"fmt"
	"os"
	"portal/app"

	"github.com/fbs-io/core"
)

func main() {
	c, err := core.New()
	if err != nil {
		fmt.Println("初始化失败, 错误:", err)
		os.Exit(2)
	}
	app.New(c)
	c.Run()

}
