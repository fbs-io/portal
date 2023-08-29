/*
 * @Author: reel
 * @Date: 2023-08-24 06:03:52
 * @LastEditors: reel
 * @LastEditTime: 2023-08-29 22:26:05
 * @Description: User表相关接口
 */
package auth

import (
	"gorm.io/gorm"
)

// 修改密码
func (user *User) chpwd(param *userChPwdParams) (err error) {
	err = user.CheckPwd(param.OldPwd)
	if err != nil {
		return
	}
	user.Password = param.NewPwd
	// pb, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	// if err != nil {
	// 	return
	// }
	// user.Password = string(pb)
	// err = tx.Save(user).Error
	// if err != nil {
	// 	return
	// }
	return
}

func (user *User) update(tx *gorm.DB, param *userUpdateParams) (err error) {
	user.NickName = param.NickName
	user.Email = param.Email
	user.Role = param.Role
	user.ID = param.ID
	user.Status = param.Status
	return tx.Updates(user).Error
}
