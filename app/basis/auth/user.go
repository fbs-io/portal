/*
 * @Author: reel
 * @Date: 2023-07-18 06:41:14
 * @LastEditors: reel
 * @LastEditTime: 2023-08-29 22:24:16
 * @Description: 用户表,管理用户信息
 */
package auth

import (
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
	UUID        string           `gorm:"comment:uuid"`
	Permissions map[string]bool  `gorm:"-" json:"permission"` // 权限校验
	rdb.Model
}

type UserList struct {
	ID       uint             `json:"id"`
	Account  string           `json:"account"`
	NickName string           `json:"nick_name"`
	Email    string           `json:"email"`
	IP       string           `json:"ip"`
	CreateAt uint64           `json:"create_at"`
	Status   int8             `json:"status"`
	Role     rdb.ModeListJson `json:"role" gorm:"type:varchar(1000)"`
}

func (ul *UserList) AfterFind(tx *gorm.DB) (err error) {

	return
}

func (u *User) TableName() string {
	return "e_auth_user"
}

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

func (u *User) CheckPwd(pwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pwd))
}

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
	}
}
