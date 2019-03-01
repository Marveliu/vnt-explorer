package models

import (
	"github.com/astaxie/beego/orm"
)

type Node struct {
	Address string `orm:"pk"`
	Vname   string `orm:"unique"`
	Home    string
	Logo    string
	Ip      string
	Status  int `orm:"index"`
	Votes   string
	Block   []*Block `orm:"reverse(many)"`
}

func (n *Node) Insert() error {
	o := orm.NewOrm()
	_, err := o.Insert(n)
	return err
}

func (n *Node) List(order string, offset, limit int, fields []string) ([]*Node, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(n)

	if order == "asc" {
		qs = qs.OrderBy("Votes")
	} else {
		qs = qs.OrderBy("-Votes")
	}

	var nodes []*Node
	_, err := qs.Offset(offset).Limit(limit).All(&nodes, fields...)
	return nodes, err
}

func (n *Node) Get(address string) (*Node, error) {
	o := orm.NewOrm()
	n.Address = address
	err := o.Read(n)
	return n, err
}
