/*
 * @Author: reel
 * @Date: 2023-08-19 17:38:01
 * @LastEditors: reel
 * @LastEditTime: 2023-08-22 06:53:51
 * @Description: 用户信息相关接口
 */
package auth

import (
	"fmt"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/fbs-io/core/store/rdb"
	"golang.org/x/crypto/bcrypt"
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
		user := GetUser(param.ID)
		tx := ctx.TX()
		if user == nil {
			user = &User{}
			err := tx.Model(user).Find(user).Error
			if err != nil {
				ctx.JSON(errno.ERRNO_RDB_QUERY.ToMapWithError(err))
				return
			}
		}
		user.NickName = param.NickName
		user.Email = param.Email
		user.Role = param.Role
		user.ID = param.ID
		err := tx.Updates(user).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_UPDATE.ToMapWithError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.ToMapWithData(user.UserInfo()))
	}
}

type userChPwdParams struct {
	ID      uint   `json:"id"  binding:"required"`
	OldPwd  string `json:"old_pwd" binding:"required" conditions:"-"`
	NewPwd  string `json:"new_pwd" binding:"required" conditions:"-"`
	NewPwd2 string `json:"new_pwd2" binding:"eqfield=NewPwd" conditions:"-"`
}

func chpwd() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*userChPwdParams)
		user := GetUser(param.ID)
		tx := ctx.TX()
		if user == nil {
			user = &User{}
			err := tx.Model(user).Find(user).Error
			if err != nil {
				ctx.JSON(errno.ERRNO_RDB_QUERY.ToMapWithError(err))
				return
			}
		}
		fmt.Println(user)
		err := user.CheckPwd(param.OldPwd)
		if err != nil {
			ctx.JSON(errno.ERRNO_AUTH_USER_OR_PWD.ToMapWithError(err))
			return
		}
		user.Password = param.NewPwd
		pb, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(errno.ERRNO_SYSTEM.ToMapWithError(err))
			return
		}
		user.Password = string(pb)
		err = tx.Save(user).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_UPDATE.ToMapWithError(err))
			return
		}

		ctx.JSON(errno.ERRNO_OK.ToMap())
	}
}
