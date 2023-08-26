/*
 * @Author: reel
 * @Date: 2023-08-19 17:38:01
 * @LastEditors: reel
 * @LastEditTime: 2023-08-26 21:31:50
 * @Description: 用户信息相关接口
 */
package auth

import (
	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/fbs-io/core/store/rdb"
)

// TODO:补充其他用户信息, 如部门等
// 设计思路: 用户和员工分开, 用户可以绑定员工, 但员工不一定有登陆账户
type userUpdateParams struct {
	ID       uint             `json:"id"  binding:"required"`
	NickName string           `json:"nick_name" conditions:"-"`
	Email    string           `json:"email" binding:"omitempty,email" conditions:"-"`
	Role     rdb.ModeListJson `json:"role" gorm:"type:varchar(1000)" conditions:"-"`
}

func userUpdate() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*userUpdateParams)
		tx := ctx.TX()
		user := GetUser(param.ID, tx)
		if user == nil {
			ctx.JSON(errno.ERRNO_RDB_QUERY)
			return
		}
		err := user.update(tx, param)
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_UPDATE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.WrapData(user.UserInfo()).Notify())
	}
}

// 密码修改参数
type userChPwdParams struct {
	ID      uint   `json:"id"  binding:"required"`
	OldPwd  string `json:"old_pwd" binding:"required" conditions:"-"`
	NewPwd  string `json:"new_pwd" binding:"required" conditions:"-"`
	NewPwd2 string `json:"new_pwd2" binding:"eqfield=NewPwd" conditions:"-"`
}

func chpwd() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*userChPwdParams)
		tx := ctx.TX()
		user := GetUser(param.ID, tx)
		if user == nil {
			ctx.JSON(errno.ERRNO_RDB_QUERY.ToMap())
			return
		}
		err := user.chpwd(tx, param)
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_UPDATE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}

// orders, page_num, page_size 作为保留字段用于条件生成
type userListQueryParams struct {
	PageNum  int    `form:"page_num"`
	PageSize int    `form:"page_size"`
	Orders   string `form:"orders"`
	Name     string `form:"nick_name" conditions:"like"`
}

func userListQuery() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*userListQueryParams)
		tx := ctx.TX()
		users := make([]*UserList, 0, 100)
		err := tx.Model(&User{}).Find(&users).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_QUERY.WrapError(err))
			return
		}
		var count int64
		tx.Model(&User{}).Count(&count)
		data := map[string]interface{}{
			"page_num":  param.PageNum,
			"page_size": param.PageSize,
			"total":     count,
			"rows":      users,
		}
		ctx.JSON(errno.ERRNO_OK.WrapData(data))
	}
}
