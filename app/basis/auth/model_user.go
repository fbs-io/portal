/*
 * @Author: reel
 * @Date: 2023-07-18 06:41:14
 * @LastEditors: reel
 * @LastEditTime: 2024-03-23 05:46:34
 * @Description: 用户表,管理用户信息
 */
package auth

import (
	"portal/pkg/consts"
	"sync"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/store/rdb"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 设计思路: 用户和员工分开, 用户可以绑定员工, 但员工不一定有登陆账户
type User struct {
	Account     string                     `gorm:"comment:账户;unique"`
	Password    string                     `gorm:"comment:密码"`
	NickName    string                     `gorm:"comment:账户名"`
	Email       string                     `gorm:"comment:邮箱"`
	IP          string                     `gorm:"comment:登陆IP"`
	Super       string                     `gorm:"comment:是否超管, Y表示是, N表示否;default:N"`
	Company     []string                   `json:"company" gorm:"-"`
	Roles       map[string][]string        `json:"-" gorm:"-"`
	Permissions map[string]map[string]bool `json:"permission" gorm:"-"`  // 权限校验, 按分区构建
	Menu        map[string][]*core.Sources `json:"menu" gorm:"-"`        // 菜单
	ManageMenu  map[string][]*core.Sources `json:"manage_menu" gorm:"-"` // 菜单
	Position1   map[string]string          `json:"position1" gorm:"-"`   // 主岗
	Position2   map[string][]string        `json:"position2" gorm:"-"`   // 兼岗
	UUID        string                     `gorm:"comment:uuid"`
	lock        *sync.RWMutex              `json:"-" gorm:"-"`
	rdb.Model
}

type UserList struct {
	ID         uint     `json:"id"`
	Account    string   `json:"account"`
	NickName   string   `json:"nick_name"`
	Email      string   `json:"email"`
	IP         string   `json:"ip"`
	Super      string   `json:"super"`
	CreatedAt  uint64   `json:"created_at"`
	Status     int8     `json:"status"`
	Position1  string   `json:"position1" gorm:"-"` //主岗
	Position2  []string `json:"position2" gorm:"-"` //兼岗
	Positions  []string `json:"positions" gorm:"-"` //所有岗位
	Department []string `json:"department" gorm:"-"`
	Company    []string `json:"company" gorm:"-"`
	Role       []string `json:"role" gorm:"-"`
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
	u.lock = &sync.RWMutex{}
	u.Company = make([]string, 0, 10)
	u.Position1 = make(map[string]string, 100)
	u.Position2 = make(map[string][]string, 100)
	u.Roles = make(map[string][]string, 10)
	u.Menu = make(map[string][]*core.Sources, 10)
	u.ManageMenu = make(map[string][]*core.Sources, 10)
	u.Permissions = make(map[string]map[string]bool, 10)
	return nil
}

func (u *User) UpdateCompany(companys ...string) {
	u.lock.Lock()
	defer u.lock.Unlock()
	if u.Company == nil {
		u.Company = make([]string, 0, 10)
	}
	u.Company = append(u.Company, companys...)
}

func (u *User) UpdateRole(sk string, role ...string) {
	u.lock.Lock()
	defer u.lock.Unlock()
	if u.Roles == nil {
		u.Roles = make(map[string][]string, 10)
	}
	if u.Roles[sk] == nil {
		u.Roles[sk] = make([]string, 0, 10)
	}
	u.Roles[sk] = append(u.Roles[sk], role...)
}

func (u *User) UpdatePosition(sk string, positions ...*RlatUserPosition) {
	u.lock.Lock()
	defer u.lock.Unlock()
	for _, position := range positions {
		if position.Account != u.Account {
			continue
		}
		if position.IsPosition == 1 {
			if u.Position1 == nil {
				u.Position1 = make(map[string]string, 10)
			}
			u.Position1[sk] = position.PositionCode
		} else {
			if u.Position2 == nil {
				u.Position2 = make(map[string][]string, 10)
			}
			if u.Position2[sk] == nil {
				u.Position2[sk] = make([]string, 0, 10)
			}
			u.Position2[sk] = append(u.Position2[sk], position.PositionCode)
		}
	}
}

func (u *User) DeleteCompany() {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.Company = nil
}

func (u *User) DeleteRoles(sk string) {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.Roles[sk] = nil
}

func (u *User) DeletePosition(sk string) {
	u.lock.Lock()
	defer u.lock.Unlock()
	delete(u.Position1, sk)
	u.Position2[sk] = make([]string, 0)
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

// 用户操作方法
// 密码修改参数
type userChPwdParams struct {
	OldPwd  string `json:"old_pwd" binding:"required" conditions:"-"`
	NewPwd  string `json:"new_pwd" binding:"required" conditions:"-"`
	NewPwd2 string `json:"new_pwd2" binding:"eqfield=NewPwd" conditions:"-"`
}

// 用户个人信息修改
type userUpdateParams struct {
	NickName string           `json:"nick_name" conditions:"-"`
	Email    string           `json:"email" binding:"omitempty,email" conditions:"-"`
	Super    string           `json:"super" conditions:"-"`
	Status   int8             `json:"status" conditions:"-"`
	Role     rdb.ModeListJson `json:"role" gorm:"type:varchar(1000)" conditions:"-"`
}

// 用户管理操作方法
// 设计思路: 用户和员工分开, 用户可以绑定员工, 但员工不一定有登陆账户
type userAddParams struct {
	Account   string   `json:"account"`
	Password  string   `json:"password"`
	NickName  string   `json:"nick_name"`
	Email     string   `json:"email" binding:"omitempty,email"`
	Super     string   `json:"super"`
	Position1 string   `json:"position1" conditions:"-"` // 主岗
	Position2 []string `json:"position2" conditions:"-"` // 兼岗
	Company   []string `json:"company"`
	Role      []string `json:"role"`
}

type usersEditParams struct {
	ID        []uint   `json:"id"  binding:"required" conditions:"-"`
	NickName  string   `json:"nick_name" conditions:"-"`
	Pwd       string   `json:"password" conditions:"-"`
	Super     string   `json:"super" conditions:"-"`
	Status    int8     `json:"status" conditions:"-"`
	Email     string   `json:"email" binding:"omitempty,email" conditions:"-"`
	Role      []string `json:"role" conditions:"-"`
	Company   []string `json:"company" conditions:"-"`
	Position1 string   `json:"position1" conditions:"-"` // 主岗
	Position2 []string `json:"position2" conditions:"-"` // 兼岗
}

// orders, page_num, page_size 作为保留字段用于条件生成
type usersQueryParams struct {
	PageNum  int    `form:"page_num"`
	PageSize int    `form:"page_size"`
	Orders   string `form:"orders"`
	Super    string `form:"super"`
	Name     string `form:"nick_name" conditions:"like"`
}
