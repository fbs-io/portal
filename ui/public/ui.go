/*
 * @Author: reel
 * @Date: 2023-10-05 00:17:15
 * @LastEditors: reel
 * @LastEditTime: 2023-10-05 06:23:38
 * @Description: 请填写简介
 */
package dist

import "embed"

var (

	//go:embed website/*
	Static embed.FS

	//go:embed index.html
	Index []byte
)
