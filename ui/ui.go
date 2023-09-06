/*
 * @Author: reel
 * @Date: 2023-09-05 20:10:17
 * @LastEditors: reel
 * @LastEditTime: 2023-09-05 23:13:13
 * @Description: 请填写简介
 */
package ui

import "embed"

var (

	//go:embed website/*
	Static embed.FS

	//go:embed website/index.html
	Index []byte
)
