/*
 * @Author: reel
 * @Date: 2023-08-19 17:38:01
 * @LastEditors: reel
 * @LastEditTime: 2023-09-05 19:00:15
 * @Description: 用户信息相关接口
 */
package auth

import (
	"time"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/fbs-io/core/store/rdb"
)

// 用户操作方法
// 密码修改参数
type userChPwdParams struct {
	OldPwd  string `json:"old_pwd" binding:"required" conditions:"-"`
	NewPwd  string `json:"new_pwd" binding:"required" conditions:"-"`
	NewPwd2 string `json:"new_pwd2" binding:"eqfield=NewPwd" conditions:"-"`
}

func chpwd() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*userChPwdParams)
		tx := ctx.TX()
		user := &User{
			Account: ctx.Auth(),
		}

		err := tx.Model(user).Where("account = (?)", user.Account).Find(user).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_QUERY)
			return
		}
		err = user.chpwd(param)
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

// 用户个人信息修改
type userUpdateParams struct {
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

		user := &User{}
		user.NickName = param.NickName
		user.Email = param.Email
		user.Role = param.Role
		user.Account = ctx.Auth()
		user.Status = param.Status
		err := tx.Where("account = (?)", ctx.Auth()).Updates(user).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_UPDATE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.WrapData(user.UserInfo()).Notify())
	}
}

// 用户管理操作方法
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
			Role:     param.Role,
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

// type usersUpdateParams struct {

// 	NickName string           `json:"nick_name" conditions:"-"`
// 	Email    string           `json:"email" binding:"omitempty,email" conditions:"-"`
// 	Super    string           `json:"super" conditions:"-"`
// 	Status   int8             `json:"status" conditions:"-"`
// 	Role     rdb.ModeListJson `json:"role" gorm:"type:varchar(1000)" conditions:"-"`
// }

// func userUpdate() core.HandlerFunc {
// 	return func(ctx core.Context) {
// 		param := ctx.CtxGetParams().(*userUpdateParams)
// 		tx := ctx.TX()

// 		user := &User{}
// 		user.NickName = param.NickName
// 		user.Email = param.Email
// 		user.Role = param.Role
// 		user.Account = ctx.Auth()
// 		user.Status = param.Status
// 		err := tx.Where("account = (?)", ctx.Auth()).Updates(user).Error
// 		if err != nil {
// 			ctx.JSON(errno.ERRNO_RDB_UPDATE.WrapError(err))
// 			return
// 		}
// 		ctx.JSON(errno.ERRNO_OK.WrapData(user.UserInfo()).Notify())
// 	}
// }

// orders, page_num, page_size 作为保留字段用于条件生成
type usersQueryParams struct {
	PageNum  int    `form:"page_num"`
	PageSize int    `form:"page_size"`
	Orders   string `form:"orders"`
	Super    string `form:"super"`
	Name     string `form:"nick_name" conditions:"like"`
}

func usersQuery() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*usersQueryParams)
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

		// for _, user := range users {
		// 	roleList := make([]interface{}, 0, 100)
		// 	for _, role := range user.Role {
		// 		r := GetRole(uint(role.(float64)), ctx.Core().RDB().DB())
		// 		roleList = append(roleList, r.Label)
		// 	}
		// 	user.Role = roleList
		// }
		ctx.JSON(errno.ERRNO_OK.WrapData(data))
	}
}

type usersUpdateParams struct {
	ID       []uint           `json:"id"  binding:"required" conditions:"-"`
	NickName string           `json:"nick_name" conditions:"-"`
	Pwd      string           `json:"password" conditions:"-"`
	Super    string           `json:"super" conditions:"-"`
	Status   int8             `json:"status" conditions:"-"`
	Email    string           `json:"email" binding:"omitempty,email"`
	Role     rdb.ModeListJson `json:"role" gorm:"type:varchar(1000)" conditions:"-"`
}

func usersUpdate() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*usersUpdateParams)
		tx := ctx.TX()
		user := &User{
			Role:     param.Role,
			Super:    param.Super,
			Password: param.Pwd,
			NickName: param.NickName,
		}
		user.Status = param.Status

		err := tx.Model(user).Where("id in (?)", param.ID).Updates(user).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_UPDATE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.WrapData(user.UserInfo()).Notify())
	}
}

type usersDeleteParams struct {
	ID []uint `json:"id" conditions:"-"`
}

// 软删除
func usersDelete() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*usersDeleteParams)
		tx := ctx.TX()

		user := &User{}
		user.DeletedBy = ctx.Auth()
		user.DeletedAT = uint(time.Now().Unix())

		err := tx.Model(user).Where("id in (?)", param.ID).Updates(user).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_DELETE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}
