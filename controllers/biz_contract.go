package controllers

import (
	"github.com/astaxie/beego"
	"github.com/vntchain/vnt-explorer/common"
	"github.com/vntchain/vnt-explorer/models"
)

type BizContractController struct {
	BaseController
}

func (c *BizContractController) List() {
	offset, err := c.GetInt64("offset")
	if err != nil {
		beego.Warn("Failed to read offset: ", err.Error())
		offset = common.DefaultOffset
	}
	limit, err := c.GetInt64("limit")
	if err != nil {
		beego.Warn("Failed to read limit: ", err.Error())
		limit = common.DefaultPageSize
	}
	order := c.GetString("order")
	fields := c.getFields()
	model := &models.BizContract{}
	ds, err := model.List(offset, limit, order, fields...)
	if err != nil {
		c.ReturnErrorMsg("Failed to list BizContracts: %s", err.Error(), "")
	} else {
		count := make(map[string]int64)
		count["count"], err = model.Count()
		if err != nil {
			c.ReturnErrorMsg("Failed to list BizContracts: %s", err.Error(), "")
			return
		}
		c.ReturnData(ds, count)
	}

}

func (c *BizContractController) Get() {
	addr := c.Ctx.Input.Param(":addr")

	fields := c.getFields()
	beego.Info("Will read colums: ", fields, "number", addr)

	model := &models.BizContract{}
	ret, err := model.Get(addr)
	if err != nil {
		c.ReturnErrorMsg("Failed to read BizContract: %s", err.Error(), "")
	} else {
		c.ReturnData(ret, nil)
	}
}

func (c *BizContractController) Count() {
	model := &models.BizContract{}
	count, err := model.Count()
	if err != nil {
		c.ReturnErrorMsg("Failed to get BizContract count: %s", err.Error(), "")
	} else {
		c.ReturnData(count, nil)
	}
}
