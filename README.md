<!--
 * @Author: reel
 * @Date: 2023-09-01 22:37:45
 * @LastEditors: reel
 * @LastEditTime: 2023-09-01 22:37:47
 * @Description: 请填写简介
-->

# FBS 企业级应用开发框架

    后端基于fbs-io/core, 前端基于scui的企业级应用开发框架, 可以快速完成业务逻辑及前端应用的开发, 同时借助core内置的资源表, 快速生成系统菜单及权限配置.

### 模块

* [X] 个人中心
* [X] 用户管理
* [X] 角色管理
* [X] 菜单管理
* [X] 权限控制
* [ ] 首页画面配置
* [ ] 全局搜索
* [ ] 字典管理
* [ ] 应用管理
* [ ] 系统日志
* [ ] 定时任务
* [ ] 表格模板

### 特性

* 包含fbs-io/core的全部特性, 内置的服务器管理功能可以直接访问
* 应用独立的web端口, 和core的管理端口进行区分.
* 完善的基于角色的权限控制, 实现数据可控
* 实现快速开发, 定义好模型, 1小时完成一张业务表格的开发, 包含CURD全部功能

```go
package main

import (
	"fmt"
	"os"

	"github.com/fbs-io/core"
)

func main() {
    c, err := core.New()
    if err != nil {
        fmt.Println("初始化失败, 错误:", err)
        os.Exit(2)
    }

    ajax := c.Group("ajax")

    dim := ajax.Group("dim", "字典数据")

    pkl := dim.Group("picklist", "码值表")

    pkl.GET("list", "获取码值列表", params{}, func(ctx core.Context) {
        ctx.JSON(errno.ERRNO_OK.ToMapWithData("请求成功"))

    })

    c.Run()

}
```
