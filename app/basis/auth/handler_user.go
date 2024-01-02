/*
 * @Author: reel
 * @Date: 2023-08-19 17:38:01
 * @LastEditors: reel
 * @LastEditTime: 2023-12-30 22:42:18
 * @Description: 用户信息相关接口
 */
package auth

import (
	"portal/app/basis/org"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/fbs-io/core/pkg/errorx"
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
	Account   string           `json:"account"`
	Password  string           `json:"password"`
	NickName  string           `json:"nick_name"`
	Email     string           `json:"email" binding:"omitempty,email"`
	Super     string           `json:"super"`
	Position1 string           `json:"position1" conditions:"-"` // 主岗
	Position2 []string         `json:"position2" conditions:"-"` // 兼岗
	Company   rdb.ModeListJson `json:"company"`
	Role      rdb.ModeListJson `json:"role"`
}

func userAdd() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*userAddParams)
		tx := ctx.TX()
		company_code := ""
		ci, ok := ctx.Ctx().Copy().Get(core.CTX_SHARDING_KEY)
		if ok && ci != nil {
			company_code = ci.(string)
		}
		rups := make([]*RlatUserPosition, 0, 10)
		user := &User{
			Account:  param.Account,
			Password: param.Password,
			NickName: param.NickName,
			Email:    param.Email,
			Super:    param.Super,
			Company:  param.Company,
			Roles:    map[string]interface{}{company_code: param.Role},

			// Departments: map[string]interface{}{company_code: param.Department},
		}
		// 主岗
		if param.Position1 != "" {
			rups = append(rups, &RlatUserPosition{
				Account:      user.Account,
				PositionCode: param.Position1,
				IsPosition:   1,
			})
		}
		//兼岗
		for _, posit := range param.Position2 {
			rup := &RlatUserPosition{
				Account:      user.Account,
				PositionCode: posit,
				IsPosition:   -1,
			}
			rups = append(rups, rup)
		}
		tx.Begin()
		err := tx.Create(user).Error
		if err != nil {
			if rdb.IsUniqueError(err) {
				ctx.JSON(errno.ERRNO_RDB_DUPLICATED_KEY.WrapError(err))
				return
			}
			ctx.JSON(errno.ERRNO_RDB_CREATE.WrapError(err))
			return
		}
		// 如果有更新, 则按用户全部更新
		if len(rups) > 0 {
			tx = ctx.NewTX()
			tx.Begin()
			err = tx.Model(&RlatUserPosition{}).Where("account = ?", param.Account).Delete(&RlatUserPosition{}).Error
			if err != nil {
				ctx.JSON(errno.ERRNO_RDB_UPDATE.WrapError(errorx.Wrap(err, "更新岗位失败")).Notify())
				return
			}
			err = tx.Model(&RlatUserPosition{}).Create(rups).Error
			if err != nil {
				ctx.JSON(errno.ERRNO_RDB_UPDATE.WrapError(errorx.Wrap(err, "批量更新岗位失败")).Notify())
				return
			}
		}
		ctx.JSON(errno.ERRNO_OK.Notify())

	}
}

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

		err := tx.Model(um).Order("id").Find(&users).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_QUERY.WrapError(err))
			return
		}

		// TODO: 判断是否需要通过数据库视图进行处理关系表逻辑, 通过程序有些复杂
		// 处理关系表,
		user_account := make([]string, 0, 100)
		user_map := make(map[string]*UserList, 100)
		for _, u := range users {
			user_account = append(user_account, u.Account)
			user_map[u.Account] = u
		}
		rlat := make([]*RlatUserPosition, 0, 10)
		ctx.NewTX().Where("account in ? ", user_account).Find(&rlat)

		position_user := make(map[string][]string)
		position_codes := make([]string, 0, 100)
		for _, r := range rlat {
			if position_user[r.PositionCode] == nil {
				position_user[r.PositionCode] = make([]string, 0, 100)
			}
			position_user[r.PositionCode] = append(position_user[r.PositionCode], r.Account)
			if (user_map[r.Account]).Position == nil {
				user_map[r.Account].Position = make(rdb.ModeListJson, 0, 100)
			}
			user_map[r.Account].Position = append(user_map[r.Account].Position, r.PositionCode)
			position_codes = append(position_codes, r.PositionCode)
			if r.IsPosition == 1 {
				user_map[r.Account].Position1 = r.PositionCode
			} else {
				if user_map[r.Account].Position2 == nil {
					user_map[r.Account].Position2 = make(rdb.ModeListJson, 0, 10)
				}
				user_map[r.Account].Position2 = append(user_map[r.Account].Position2, r.PositionCode)
			}

		}

		// 获取岗位表
		positions := make([]*org.Position, 0, 100)

		ctx.NewTX().Where("position_code in ? ", position_codes).Find(&positions)

		for _, position := range positions {
			for _, account := range position_user[position.PositionCode] {
				if user_map[account].Department == nil {
					user_map[account].Department = make(rdb.ModeListJson, 0, 100)
				}
				user_map[account].Department = append(user_map[account].Department, position.DepartmentCode)
			}
		}

		// 获取总数
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

type usersUpdateParams struct {
	ID        []uint           `json:"id"  binding:"required" conditions:"-"`
	NickName  string           `json:"nick_name" conditions:"-"`
	Pwd       string           `json:"password" conditions:"-"`
	Super     string           `json:"super" conditions:"-"`
	Status    int8             `json:"status" conditions:"-"`
	Email     string           `json:"email" binding:"omitempty,email" conditions:"-"`
	Position1 string           `json:"position1" conditions:"-"` // 主岗
	Position2 []string         `json:"position2" conditions:"-"` // 兼岗
	Role      rdb.ModeListJson `json:"role" conditions:"-"`
	Company   rdb.ModeListJson `json:"company" conditions:"-"`
}

func usersUpdate() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*usersUpdateParams)
		tx := ctx.TX()
		users := make([]*User, 0, 100)
		err := tx.Model(&User{}).Where("id in (?)", param.ID).Find(&users).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_UPDATE.WrapError(err))
			return
		}
		company_code := ""
		ci, ok := ctx.Ctx().Copy().Get(core.CTX_SHARDING_KEY)
		if ok && ci != nil {
			company_code = ci.(string)
		}
		rups := make([]*RlatUserPosition, 0, 10)
		// 删除用
		dus := make([]string, 0, 10)
		for _, user := range users {
			dus = append(dus, user.Account)
			user.Super = param.Super
			user.NickName = param.NickName
			user.Company = param.Company
			user.Email = param.Email
			user.Status = param.Status
			user.Password = param.Pwd
			user.Roles[company_code] = param.Role
			if param.Position1 != "" {
				rups = append(rups, &RlatUserPosition{
					Account:      user.Account,
					PositionCode: param.Position1,
					IsPosition:   1,
				})
			}
			//兼岗
			for _, posit := range param.Position2 {
				rup := &RlatUserPosition{
					Account:      user.Account,
					PositionCode: posit,
					IsPosition:   -1,
				}
				rups = append(rups, rup)
			}
			// user.Departments[company_code] = param.Department
			err := tx.Model(user).Where("id = (?)", user.ID).Updates(user).Error
			if err != nil {
				ctx.JSON(errno.ERRNO_RDB_UPDATE.WrapError(err))
				return
			}
		}
		// 如果有更新, 则按用户全部更新
		if len(rups) > 0 {
			tx = ctx.NewTX()
			tx.Begin()
			err = tx.Model(&RlatUserPosition{}).Where("account in ? ", dus).Delete(&RlatUserPosition{}).Error
			if err != nil {
				ctx.JSON(errno.ERRNO_RDB_UPDATE.WrapError(errorx.Wrap(err, "更新岗位失败")).Notify())
				return
			}
			err = tx.Model(&RlatUserPosition{}).Create(rups).Error
			if err != nil {
				ctx.JSON(errno.ERRNO_RDB_UPDATE.WrapError(errorx.Wrap(err, "批量更新岗位失败")).Notify())
				return
			}
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
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
		err := tx.Model(user).Where("id in (?)", param.ID).Delete(user).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_DELETE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}
