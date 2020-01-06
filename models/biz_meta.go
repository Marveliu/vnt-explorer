package models

import (
	"github.com/astaxie/beego/orm"
)

type BizMeta struct {
	No        uint32 `orm:"pk"`
	BizName   string
	BizType   string
	Desc      string
	Datas     string `orm:"type(text)"`
	Tasks     string `orm:"type(text)"`
	Timestamp uint64
}

func (b *BizMeta) Insert() error {
	o := orm.NewOrm()
	_, err := o.InsertOrUpdate(b)
	return err
}

func (b *BizMeta) Get(No int) (*BizMeta, error) {
	o := orm.NewOrm()
	b.No = uint32(No)
	err := o.Read(b)
	return b, err
}

func (b *BizMeta) List(offset, limit int64, order string, fields ...string) ([]*BizMeta, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(b)
	if order == "asc" {
		qs = qs.OrderBy("Timestamp")
	} else {
		qs = qs.OrderBy("-Timestamp")
	}
	var bizMetas []*BizMeta
	_, err := qs.Offset(offset).Limit(limit).All(&bizMetas, fields...)
	return bizMetas, err
}

func (b *BizMeta) Count() (int64, error) {
	o := orm.NewOrm()
	cnt, err := o.QueryTable(b).Count()
	return cnt, err
}
