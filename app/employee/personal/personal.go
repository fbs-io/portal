/*
 * @Author: reel
 * @Date: 2023-06-24 14:46:01
 * @LastEditors: reel
 * @LastEditTime: 2023-07-18 06:57:34
 * @Description: 员工个人信息管理
 */
package personal

import (
	"github.com/fbs-io/core/store/rdb"
)

type PersonalBase struct {
	FirstName string `gorm:"firstname"` // 姓
	LastName  string `gorm:"lastname"`  // 名
	FullName  string `gorm:"fullname"`  // 姓名
	Gender    string `gorm:"gender"`    // 性别
	UUID      string `gorm:"uuid"`      // 唯一编码
}

type Personal struct {
	PersonalBase
	rdb.Model
}

func (p *Personal) TableName() string {
	return "e_employee_personal"
}
