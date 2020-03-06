package models

import (
	"github.com/astaxie/beego"
	"testing"
)

func init() {
	beego.InitBeegoBeforeTest("/Users/mac/gopath/src/github.com/vntchain/vnt-explorer/conf/app.conf")
}

func TestDetect_Insert(t *testing.T) {
	m := Detect{
		Addr:      "12341234",
		Score:     0.96,
		Detail:    "你错了",
		Type:      1,
		TimeStamp: 12341234,
	}
	m.Insert()
}
