package models

import (
	"github.com/astaxie/beego/orm"
)

type BizMeta struct {
	No        uint32 `orm:"pk"`
	BizName   string
	BizType   string
	Desc      string
	Datas     string
	Tasks     string
	Timestamp uint64
}

func (t *BizMeta) Insert() error {
	o := orm.NewOrm()
	_, err := o.InsertOrUpdate(t)
	return err
}
