package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
)

type ErrorBlockNumber struct {
	format string
	number string
}

func (e ErrorBlockNumber) Error() string {
	return fmt.Sprintf(e.format, e.number)
}

type Block struct {
	Number         uint64 `orm:"pk"`
	TimeStamp      uint64
	TxCount        int
	Hash           string `orm:"unique"`
	ParentHash     string
	Producer       string `orm:"index"`
	ProducerDetail *Node  `orm:"-"`
	Size           string
	GasUsed        uint64
	GasLimit       uint64
	BlockReward    string
	Reward		   float64
	Fee			   float64
	ExtraData      string
	Tps            float32
	Witnesses      []*Node `orm:"rel(m2m)"`
}

func (b *Block) Insert() error {
	o := orm.NewOrm()
	_, err := o.InsertOrUpdate(b)
	return err
}

func (b *Block) Update() error {
	o := orm.NewOrm()
	_, err := o.Update(b)
	return err
}

func (b *Block) List(offset, limit int64, order string, fields ...string) ([]*Block, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(b)

	if order == "asc" {
		qs = qs.OrderBy("Number")
	} else {
		qs = qs.OrderBy("-Number")
	}

	var blocks []*Block
	_, err := qs.Offset(offset).Limit(limit).All(&blocks, fields...)
	for _, b := range blocks {
		if b.Producer != "" {
			b.ProducerDetail = &Node{
				Address: b.Producer,
			}
			o.Read(b.ProducerDetail)
		}
	}
	return blocks, err
}

func (b *Block) Get(nOrh string, fields ...string) (*Block, error) {
	o := orm.NewOrm()

	var err error
	if strings.HasPrefix(nOrh, "0x") {
		beego.Info("Will read block by hash: ", nOrh)
		b.Hash = nOrh
		err = o.Read(b, "Hash")
	} else {
		beego.Info("Will read block by number: ", nOrh)
		b.Number, err = strconv.ParseUint(nOrh, 10, 64)
		if err != nil {
			e := ErrorBlockNumber{"Wrong block number: %s", nOrh}
			beego.Error(e.Error())
			return nil, e
		}
		err = o.Read(b, "Number")
	}
	if err == nil && b.Producer != "" {
		b.ProducerDetail = &Node{
			Address: b.Producer,
		}
		o.Read(b.ProducerDetail)
	}
	return b, err
}

func (b *Block) GetByNumber(number uint64) (*Block, error) {
	o := orm.NewOrm()
	b.Number = number
	err := o.Read(b, "Number")
	return b, err
}

func (b *Block) CountBellow(number uint64) (int64, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(b).SetCond(orm.NewCondition().And("number__lte", number))

	cnt, err := qs.Count()
	return cnt, err
}

func (b *Block) Last() (*Block, error) {
	o := orm.NewOrm()

	qs := o.QueryTable(b).OrderBy("-Number").Limit(1)

	var blocks []*Block
	_, err := qs.All(&blocks)

	if err != nil {
		return nil, err
	}

	if len(blocks) == 0 {
		return nil, nil
	}

	return blocks[0], nil
}

func (b *Block) TopTpsBlock() (*Block, error) {
	o := orm.NewOrm()

	qs := o.QueryTable(b).OrderBy("-tps").Limit(1)

	var blocks []*Block
	_, err := qs.All(&blocks)

	if err != nil {
		return nil, err
	}

	if len(blocks) == 0 {
		return nil, nil
	}

	return blocks[0], nil
}

func (b *Block) Count() (int64, error) {
	o := orm.NewOrm()
	//cnt, err := o.QueryTable(b).Count()
	//return cnt, err
	var list orm.ParamsList
	_, err := o.Raw("SELECT MAX(number) FROM block").ValuesFlat(&list)
	if err != nil {
		beego.Error("block table query max block num failed", err.Error())
		return 0, err
	}
	if list[0] != nil {
		num, err := strconv.Atoi(list[0].(string))
		return int64(num + 1), err
	} else {
		return 0, fmt.Errorf("block count failed")
	}
}
