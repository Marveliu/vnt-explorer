package models

import (
	"github.com/astaxie/beego"
	"testing"
)

func init() {
	beego.InitBeegoBeforeTest("/Users/mac/gopath/src/github.com/vntchain/vnt-explorer/conf/app.conf")
}

func Test_Insert(t *testing.T) {
	meta := BizMeta{
		No:        123,
		BizName:   "123",
		BizType:   "12341",
		Desc:      "1234",
		Datas:     "1231234",
		Tasks:     "1235125",
		Timestamp: 0,
	}
	meta.Insert()
}
