package models

type Report struct {
	ContractAddr string // when transaction is a contract creation
	BizType      string
	Data         string `orm:"type(text)"`
	BlockNumber  uint64
	TxHash       string
}
