package controllers

import (
	"github.com/astaxie/beego"
	"github.com/vntchain/vnt-explorer/common"
	"github.com/vntchain/vnt-explorer/models"
	"strconv"
)

type ReportController struct {
	BaseController
}

func (c *ReportController) List() {
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
	model := &models.Report{}
	ds, err := model.List(offset, limit, order, fields...)
	if err != nil {
		c.ReturnErrorMsg("Failed to list reports: %s", err.Error(), "")
	} else {
		count := make(map[string]int64)
		count["count"], err = model.Count()
		if err != nil {
			c.ReturnErrorMsg("Failed to list reports: %s", err.Error(), "")
			return
		}
		c.ReturnData(ds, count)
	}

}

func (c *ReportController) Get() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	fields := c.getFields()
	beego.Info("Will read colums: ", fields, "number", id)

	model := &models.Report{}
	ret, err := model.Get(id)
	if err != nil {
		c.ReturnErrorMsg("Failed to read report: %s", err.Error(), "")
	} else {
		c.ReturnData(ret, nil)
	}
}

func (c *ReportController) Count() {
	model := &models.Report{}
	count, err := model.Count()
	if err != nil {
		c.ReturnErrorMsg("Failed to get report count: %s", err.Error(), "")
	} else {
		c.ReturnData(count, nil)
	}
}
