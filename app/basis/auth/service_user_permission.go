/*
 * @Author: reel
 * @Date: 2024-03-21 05:44:23
 * @LastEditors: reel
 * @LastEditTime: 2024-03-22 21:54:24
 * @Description: 处理用户权限相关逻辑
 */
package auth

import (
	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errorx"
	"github.com/fbs-io/core/store/rdb"
	"gorm.io/gorm"
)

// 获取资源权限表
//
// 分为菜单管理权限, 接口权限, 设置菜单权限等
func (srv *userService) GetResourcePermission(tx *gorm.DB, account string) (menuList []*core.Sources, manageList []*core.Sources, permissions map[string]bool, err error) {
	sk := rdb.GetShardingKey(tx)
	user := srv.GetByCode(account)
	if user == nil || user.Status < 1 {
		return nil, nil, nil, errorx.Errorf("无有效用户:%s, 请检查用户是否存在或失效", account)
	}
	var permissionMap = make(map[uint]bool, 100)
	// 从角色中获取用户的资源权限列表,如果角色列表为空, 则不过滤
	if user.Super == "Y" {
		for _, item := range ResourceService.GetAllList() {
			permissionMap[item.ID] = true
		}
	} else {
		for _, roleCode := range user.Roles[sk] {
			role := RoleService.GetByCode(tx, roleCode)
			if role == nil || role.Status < 1 {
				continue
			}
			for _, codeI := range role.Sources {
				id, ok := codeI.(float64)
				if ok {
					permissionMap[uint(id)] = true
				}
			}
		}
	}
	// 从权限列表中,获取用户对应的菜单权限,按钮(接口)权限, 菜单管理权限等
	return ResourceService.GetSource(permissionMap)
}
