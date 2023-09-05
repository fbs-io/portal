/*
 * @Author: reel
 * @Date: 2023-08-20 18:30:14
 * @LastEditors: reel
 * @LastEditTime: 2023-09-05 05:52:22
 * @Description: 权限获取
 */
package auth

// // TODO: 增加角色判断
// func getPermission(c core.Core, account string) (permissionMap map[string]bool, permissionList []string, err error) {
// 	source := &core.Sources{}

// 	// TODO: 通过缓存/视图的方式, 减少数据库查询次数
// 	// 获取用户
// 	user := &User{}
// 	err = c.RDB().DB().Where("account = (?)", account).Find(user).Error
// 	if err != nil {
// 		return
// 	}

// 	// 构筑条件
// 	cond := rdb.NewCondition()
// 	cond.PageSize = 1000
// 	cond.TableName = source.TableName()
// 	cond.Where = map[string]interface{}{
// 		"is_router = ? ": core.SOURCE_ROUTER_NAN,
// 		"type in (?)":    []int8{core.SOURCE_TYPE_PERMISSION, core.SOURCE_TYPE_UNPERMISSION},
// 	}
// 	permissionList = []string{""}
// 	if user.Super != "Y" {
// 		// 获取用户关联的角色
// 		roles := make([]*Role, 0, 100)
// 		userRoles := []string{}
// 		for _, roleCode := range user.Role {
// 			userRoles = append(userRoles, roleCode.(string))
// 		}
// 		err = c.RDB().DB().Where("code in (?)", userRoles).Find(&roles).Error
// 		if err != nil {
// 			return
// 		}

// 		for _, role := range roles {
// 			for _, source := range role.Sources {
// 				permissionList = append(permissionList, source.(string))
// 			}
// 		}
// 		cond.Where["code in (?)"] = permissionList
// 	}

// 	sourceList := make([]*core.Sources, 0, 1000)
// 	err = c.RDB().BuildQuery(cond).Find(&sourceList).Error
// 	if err != nil {
// 		return
// 	}
// 	permissionMap = make(map[string]bool, 100)
// 	permissionList = make([]string, 0, 100)
// 	for _, source := range sourceList {
// 		permissionMap[source.GenRequestKey()] = true
// 		permissionList = append(permissionList, source.Code)
// 	}
// 	return
// }
