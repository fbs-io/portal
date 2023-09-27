package auth

import (
	"portal/pkg/consts"

	"github.com/fbs-io/core/store/rdb"
	"gorm.io/gorm"
)

type Company struct {
	CompanyCode        string
	CompanyName        string
	Description        string
	EffectiveStartDate int64
	EffectiveEndDate   int64
	rdb.Model
}

func (c *Company) TableName() string {
	return consts.TABLE_BASIS_AUTH_COMPANY
}

// gorm 中间件操作
func (c *Company) BeforeCreate(tx *gorm.DB) error {
	c.Model.BeforeCreate(tx)
	return nil
}

func (c *Company) BeforeUpdate(tx *gorm.DB) error {
	c.Model.BeforeUpdate(tx)
	return nil
}

func (c *Company) BeforeDelete(tx *gorm.DB) error {
	c.Model.BeforeDelete(tx)
	return nil
}
