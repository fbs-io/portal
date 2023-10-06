/*
 * @Author: reel
 * @Date: 2023-07-18 06:41:14
 * @LastEditors: reel
 * @LastEditTime: 2023-10-06 07:56:33
 * @Description: 用户表,管理用户信息
 */
package auth

import (
	"portal/pkg/consts"

	"github.com/fbs-io/core/store/rdb"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 设计思路: 用户和员工分开, 用户可以绑定员工, 但员工不一定有登陆账户
type User struct {
	Account     string           `gorm:"comment:账户;unique"`
	Password    string           `gorm:"comment:密码"`
	NickName    string           `gorm:"comment:账户名"`
	Email       string           `gorm:"comment:邮箱"`
	IP          string           `gorm:"comment:登陆IP"`
	Super       string           `gorm:"comment:是否超管, Y表示是, N表示否;default:N"`
	Role        rdb.ModeListJson `gorm:"comment:角色;type:varchar(1000)"`
	Company     rdb.ModeListJson `json:"company" gorm:"type:varchar(10240)"`
	UUID        string           `gorm:"comment:uuid"`
	Permissions map[string]bool  `gorm:"-" json:"permission"` // 权限校验
	rdb.Model
}

type UserList struct {
	ID        uint             `json:"id"`
	Account   string           `json:"account"`
	NickName  string           `json:"nick_name"`
	Email     string           `json:"email"`
	IP        string           `json:"ip"`
	Super     string           `json:"super"`
	CreatedAt uint64           `json:"created_at"`
	Status    int8             `json:"status"`
	Role      rdb.ModeListJson `json:"role" gorm:"type:varchar(10240)"`
	Company   rdb.ModeListJson `json:"company" gorm:"type:varchar(10240)"`
}

func (u *User) TableName() string {
	return consts.TABLE_BASIS_AUTH_USER

}

// gorm中间件操作

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.UUID = uuid.New().String()
	u.Model.BeforeCreate(tx)
	if u.Super == "" {
		u.Super = "N"
	}
	return u.encodePwd()
}
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.Model.BeforeUpdate(tx)
	return u.encodePwd()
}

func (u *User) BeforeDelete(tx *gorm.DB) error {
	u.Model.BeforeDelete(tx)
	return nil
}

// User模型相关操作

// 校验密码是否正确
func (u *User) CheckPwd(pwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pwd))
}

// 加密密码
func (u *User) encodePwd() error {
	if u.Password != "" {
		pb, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(pb)
		return nil
	}
	return nil
}

// 生成User用户相关信息的Map
func (u *User) UserInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":        u.ID,
		"nick_name": u.NickName,
		"account":   u.Account,
		"email":     u.Email,
		"super":     u.Super,
	}
}

// 修改密码
func (user *User) chpwd(param *userChPwdParams) (err error) {
	err = user.CheckPwd(param.OldPwd)
	if err != nil {
		return
	}
	user.Password = param.NewPwd
	return
}
