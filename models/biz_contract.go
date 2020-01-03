package models

import (
	"github.com/astaxie/beego/orm"
)

type BizContract struct {
	Address   string `orm:"pk"`
	Owner     string
	Name      string
	Desc      string
	Status    uint32
	TimeStamp uint64
	BizNo     uint32
}

func (t *BizContract) Insert() error {
	o := orm.NewOrm()
	_, err := o.InsertOrUpdate(t)
	return err
}
