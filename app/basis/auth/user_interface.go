/*
 * @Author: reel
 * @Date: 2023-08-24 06:03:52
 * @LastEditors: reel
 * @LastEditTime: 2023-08-24 06:16:55
 * @Description: User表相关接口
 */
package auth

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 修改密码
func (user *User) chpwd(tx *gorm.DB, param *userChPwdParams) (err error) {
	err = user.CheckPwd(param.OldPwd)
	if err != nil {
		return
	}
	user.Password = param.NewPwd
	pb, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	user.Password = string(pb)
	err = tx.Save(user).Error
	if err != nil {
		return
	}
	return
}

func (user *User) update(tx *gorm.DB, param *userUpdateParams) (err error) {
	user.NickName = param.NickName
	user.Email = param.Email
	user.Role = param.Role
	user.ID = param.ID
	return tx.Updates(user).Error
}
