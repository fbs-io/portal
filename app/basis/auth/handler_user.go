/*
 * @Author: reel
 * @Date: 2023-08-19 17:38:01
 * @LastEditors: reel
 * @LastEditTime: 2023-08-29 22:28:26
 * @Description: 用户信息相关接口
 */
package auth

import (
	"time"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/fbs-io/core/store/rdb"
)

// TODO:补充其他用户信息, 如部门等
// 设计思路: 用户和员工分开, 用户可以绑定员工, 但员工不一定有登陆账户
type userAddParams struct {
	Account  string           `json:"account"`
	Password string           `json:"password"`
	NickName string           `json:"nick_name"`
	Email    string           `json:"email" binding:"omitempty,email"`
	Super    string           `json:"super"`
	Role     rdb.ModeListJson `json:"role"`
}

func userAdd() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*userAddParams)
		tx := ctx.TX()
		user := &User{
			Account:  param.Account,
			Password: param.Password,
			NickName: param.NickName,
			Email:    param.Email,
			Super:    param.Super,
		}
		err := tx.Create(user).Error
		if err != nil {
			if rdb.IsUniqueError(err) {
				ctx.JSON(errno.ERRNO_RDB_DUPLICATED_KEY.WrapError(err))
				return
			}
			ctx.JSON(errno.ERRNO_RDB_CREATE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())

	}
}

type userUpdateParams struct {
	ID       uint             `json:"id"  binding:"required"`
	NickName string           `json:"nick_name" conditions:"-"`
	Email    string           `json:"email" binding:"omitempty,email" conditions:"-"`
	Super    string           `json:"super" conditions:"-"`
	Status   int8             `json:"status" conditions:"-"`
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
		err := user.chpwd(param)
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_UPDATE.WrapError(err))
			return
		}
		err = tx.Save(user).Error
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
		um := &User{}
		tx := ctx.TX(
			core.SetTxMode(core.TX_QRY_MODE_SUBID),
			core.SetTxSubTable(um.TableName()),
		)
		users := make([]*UserList, 0, 100)

		err := tx.Model(um).Find(&users).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_QUERY.WrapError(err))
			return
		}
		var count int64
		ctx.TX().Model(&User{}).Offset(-1).Limit(-1).Count(&count)
		data := map[string]interface{}{
			"page_num":  param.PageNum,
			"page_size": param.PageSize,
			"total":     count,
			"rows":      users,
		}
		ctx.JSON(errno.ERRNO_OK.WrapData(data))
	}
}

type userListUpdateParams struct {
	ID     []uint           `json:"id"  binding:"required" conditions:"-"`
	Pwd    string           `json:"password" conditions:"-"`
	Super  string           `json:"super" conditions:"-"`
	Status int8             `json:"status" conditions:"-"`
	Role   rdb.ModeListJson `json:"role" gorm:"type:varchar(1000)" conditions:"-"`
}

func usersUpdate() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*userListUpdateParams)
		tx := ctx.TX()
		user := &User{
			Role:     param.Role,
			Super:    param.Super,
			Password: param.Pwd,
		}
		user.Status = param.Status

		// if param.Pwd!=""{
		// 	user.Password =
		// }

		err := tx.Model(user).Where("id in (?)", param.ID).Updates(user).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_UPDATE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.WrapData(user.UserInfo()).Notify())
	}
}

type userDelParams struct {
	ID []uint `json:"id" conditions:"-"`
}

// 软删除
func usersDel() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*userDelParams)
		tx := ctx.TX()

		user := &User{}
		user.DeleteBy = ctx.Auth()
		user.DeletedAT = uint(time.Now().Unix())

		err := tx.Model(user).Where("id in (?)", param.ID).Updates(user).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_DELETE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}
