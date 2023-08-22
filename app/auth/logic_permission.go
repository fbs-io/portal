/*
 * @Author: reel
 * @Date: 2023-08-20 18:30:14
 * @LastEditors: reel
 * @LastEditTime: 2023-08-20 23:46:40
 * @Description: 权限获取
 */
package auth

import (
	"github.com/fbs-io/core"
	"github.com/fbs-io/core/store/rdb"
)

// TODO: 增加角色判断
func getPermission(c core.Core, u *User) (permissionMap map[string]bool, permissionList []string, err error) {
	source := &core.Sources{}

	// 构筑条件
	cond := rdb.NewCondition()
	cond.PageSize = 1000
	cond.TableName = source.TableName()
	cond.Where = map[string]interface{}{
		"is_router = ? ": core.SOURCE_ROUTER_NAN,
		"type  = ?":      core.SOURCE_TYPE_PERMISSION,
	}

	sourceList := make([]*core.Sources, 0, 1000)
	err = c.RDB().BuildQuery(cond).Find(&sourceList).Error
	if err != nil {
		return
	}
	permissionMap = make(map[string]bool, 100)
	permissionList = make([]string, 0, 100)
	for _, source := range sourceList {
		permissionMap[source.GenRequestKey()] = true
		permissionList = append(permissionList, source.Code)
	}
	return
}
