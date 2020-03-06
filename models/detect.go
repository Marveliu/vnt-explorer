package models

import (
	"github.com/astaxie/beego/orm"
)

type Detect struct {
	Id        int
	Addr      string
	Score     float64
	Detail    string
	Type      int
	TimeStamp uint64
}

func (r *Detect) Insert() error {
	o := orm.NewOrm()
	_, err := o.InsertOrUpdate(r)
	return err
}

func (r *Detect) Update() error {
	o := orm.NewOrm()
	_, err := o.Update(r)
	return err
}

func (r *Detect) Get(id int) (*Detect, error) {
	o := orm.NewOrm()
	r.Id = id
	err := o.Read(r)
	return r, err
}

func (r *Detect) List(offset, limit int64, order string, fields ...string) ([]*Detect, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(r)
	if order == "asc" {
		qs = qs.OrderBy("TimeStamp")
	} else {
		qs = qs.OrderBy("-TimeStamp")
	}
	var detects []*Detect
	_, err := qs.Offset(offset).Limit(limit).All(&detects, fields...)
	return detects, err
}

func (r *Detect) Count() (int64, error) {
	o := orm.NewOrm()
	cnt, err := o.QueryTable(r).Count()
	return cnt, err
}
