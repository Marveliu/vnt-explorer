package models

import (
	"github.com/astaxie/beego/orm"
)

type BizContract struct {
	Address   string `orm:"pk"`
	Owner     string
	Name      string
	Desc      string `orm:"type(text)"`
	Status    uint32
	TimeStamp uint64
	BizNo     uint32
}

func (b *BizContract) Insert() error {
	o := orm.NewOrm()
	_, err := o.InsertOrUpdate(b)
	return err
}

func (b *BizContract) Get(Address string) (*BizContract, error) {
	o := orm.NewOrm()
	b.Address = Address
	err := o.Read(b)
	return b, err
}

func (b *BizContract) List(offset, limit int64, order string, fields ...string) ([]*BizContract, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(b)
	if order == "asc" {
		qs = qs.OrderBy("Timestamp")
	} else {
		qs = qs.OrderBy("-Timestamp")
	}
	var bizContracts []*BizContract
	_, err := qs.Offset(offset).Limit(limit).All(&bizContracts, fields...)
	return bizContracts, err
}

func (b *BizContract) Count() (int64, error) {
	o := orm.NewOrm()
	cnt, err := o.QueryTable(b).Count()
	return cnt, err
}
