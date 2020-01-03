package models

import (
	"github.com/astaxie/beego/orm"
)

type Report struct {
	Id           int
	ContractAddr string
	MetaNo       uint64
	Data         string `orm:"type(text)"`
	BlockNumber  uint64
	TxHash       string
	TimeStamp    uint64

	BizMeta *BizMeta `orm:"-"`
}

func (r *Report) Insert() error {
	o := orm.NewOrm()
	_, err := o.InsertOrUpdate(r)
	return err
}

func (r *Report) Update() error {
	o := orm.NewOrm()
	_, err := o.Update(r)
	return err
}

func (r *Report) Get(id int) (*Report, error) {
	o := orm.NewOrm()
	r.Id = id
	err := o.Read(r)
	return r, err
}

func (r *Report) List(offset, limit int64, order string, fields ...string) ([]*Report, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(r)
	if order == "asc" {
		qs = qs.OrderBy("BlockNumber")
	} else {
		qs = qs.OrderBy("-BlockNumber")
	}
	var reports []*Report
	_, err := qs.Offset(offset).Limit(limit).All(&reports, fields...)
	for _, r := range reports {
		if r.MetaNo != 0 {
			r.BizMeta = &BizMeta{
				No: uint32(r.MetaNo),
			}
			o.Read(r.BizMeta)
		}
	}
	return reports, err
}

func (r *Report) Count() (int64, error) {
	o := orm.NewOrm()
	cnt, err := o.QueryTable(r).Count()
	return cnt, err
}
