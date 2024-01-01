/*
 * @Author: reel
 * @Date: 2023-07-18 06:41:14
 * @LastEditors: reel
 * @LastEditTime: 2023-12-31 23:01:14
 * @Description: 用户表,管理用户信息
 */
package auth

import (
	"portal/pkg/consts"

	"github.com/fbs-io/core"
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
	Company     rdb.ModeListJson `json:"company" gorm:"type:varchar(1024)"`
	Roles       rdb.ModeMapJson  `json:"-" gorm:"comment:角色;type:varchar(1000)"`
	Departments rdb.ModeMapJson  `json:"-" gorm:"type:varchar(1024)"`
	Role        rdb.ModeListJson `json:"role" gorm:"-"`
	Department  string           `json:"department" gorm:"-"`
	UUID        string           `gorm:"comment:uuid"`
	Permissions map[string]bool  `gorm:"-" json:"permission"` // 权限校验
	Menu        []*menuTree      `gorm:"-" json:"menu"`       // 菜单
	rdb.Model
}

type UserList struct {
	ID         uint             `json:"id"`
	Account    string           `json:"account"`
	NickName   string           `json:"nick_name"`
	Email      string           `json:"email"`
	IP         string           `json:"ip"`
	Super      string           `json:"super"`
	CreatedAt  uint64           `json:"created_at"`
	Status     int8             `json:"status"`
	Position1  string           `json:"position1" gorm:"-"` //主岗
	Position2  rdb.ModeListJson `json:"position2" gorm:"-"` //兼岗
	Position   rdb.ModeListJson `json:"position" gorm:"-"`  //所有岗位
	Roles      rdb.ModeMapJson  `json:"-" gorm:"type:varchar(10240)"`
	Role       rdb.ModeListJson `json:"role" gorm:"-"`
	Company    rdb.ModeListJson `json:"company" gorm:"type:varchar(10240)"`
	Department rdb.ModeListJson `json:"department" gorm:"-"`
}

func (u *User) TableName() string {
	return consts.TABLE_BASIS_AUTH_USER

}

// gorm中间件操作

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.UUID = uuid.New().String()
	if u.Super == "" {
		u.Super = "N"
	}
	return u.encodePwd()
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	return u.encodePwd()
}

func (u *User) AfterFind(tx *gorm.DB) error {
	ci, ok := tx.Get(core.CTX_SHARDING_KEY)
	if ok && ci != nil {
		if u.Roles[ci.(string)] != nil {
			u.Role = u.Roles[ci.(string)].([]interface{})
		}
	}
	if u.Departments == nil || len(u.Departments) == 0 {
		u.Departments = make(rdb.ModeMapJson, 100)
	}
	if u.Roles == nil || len(u.Roles) == 0 {
		u.Roles = make(rdb.ModeMapJson, 100)
	}
	return nil
}

func (u *UserList) AfterFind(tx *gorm.DB) error {
	ci, ok := tx.Get(core.CTX_SHARDING_KEY)

	if ok && u.Roles[ci.(string)] != nil {
		u.Role = u.Roles[ci.(string)].([]interface{})
	}

	if u.Roles == nil || len(u.Roles) == 0 {
		u.Roles = make(rdb.ModeMapJson, 100)
	}
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

// 用户和岗位关系表
// TODO: 用户和员工拆分开, 岗位和员工关联
type RlatUserPosition struct {
	Account      string `gorm:"comment:用户code;index"`
	PositionCode string `gorm:"comment:岗位code;index"`
	IsPosition   int8   `gorm:"comment:是否主岗"`
	rdb.Model
	rdb.ShardingModel
}

func (model *RlatUserPosition) TableName() string {
	return consts.TABLE_BASIS_RLAT_USER_POSITION
}
