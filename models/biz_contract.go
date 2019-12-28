package models

type BizContract struct {
	Address   string `orm:"pk"`
	Owner     string // 所有者地址
	Name      string // 合约名称
	Desc      string // 描述
	BizType   uint32 // 合约类型
	Status    uint32 // 状态
	TimeStamp uint64
}
