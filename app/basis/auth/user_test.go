/*
 * @Author: reel
 * @Date: 2023-07-18 06:58:47
 * @LastEditors: reel
 * @LastEditTime: 2023-07-20 22:43:54
 * @Description: 请填写简介
 */
package auth

import (
	"fmt"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestUser(t *testing.T) {
	p := &User{}
	p.Password = "123456"
	p.BeforeCreate(nil)
	println(p.Password)
	err := bcrypt.CompareHashAndPassword([]byte("$2a$10$gHP7RGGY5ql5e5h5IZCprOauF20RC.wFt/.VdIC7uiFpJLsQHUOGa"), []byte("123456"))
	fmt.Println(err)
	fmt.Println(p.Password)
}
