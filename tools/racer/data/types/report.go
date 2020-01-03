package types

import (
	"encoding/json"
)

type ReportField struct {
	FieldType byte
	Value     interface{}
}

type StructReport struct {
	Addr      string
	MetaNo    uint32
	BizType   string
	Datas     []ReportField
	TimeStamp uint64
}

func (s *StructReport) GetDatas() string {
	if bs, err := json.Marshal(s.Datas); err == nil {
		return string(bs)
	}
	return ""
}
