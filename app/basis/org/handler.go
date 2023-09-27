/*
 * @Author: reel
 * @Date: 2023-09-19 04:34:39
 * @LastEditors: reel
 * @LastEditTime: 2023-09-27 22:01:56
 * @Description: api操作
 */
package org

import (
	"portal/pkg/consts"
	"portal/pkg/sequence"

	"github.com/fbs-io/core"
)

type OrgAddParams struct {
	OrgCode      string
	OrgName      string
	OrgShortName string
	OrgComment   string
	OrgBusiness  string
}

func add() core.HandlerFunc {
	return func(ctx core.Context) {
		seq := sequence.New(ctx.Core(), "test", sequence.SetDateFormat(consts.DATE_FORMAT_SHORT_YMDHM), sequence.SetPrefix("C"), sequence.SetSequenceLen(5))
		code := seq.Code()
		ctx.JSON(code)
	}
}
