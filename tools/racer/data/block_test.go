package data

import (
	"github.com/astaxie/beego"
	"testing"
)

func init() {
	beego.InitBeegoBeforeTest("/Users/mac/gopath/src/github.com/vntchain/vnt-explorer/conf/app.conf")
}

func Test_getBlock(t *testing.T) {
	GetBlock(25534)
}

func Test_getTx(t *testing.T) {
	_, txs, _ := GetBlock(25534)
	for _, v := range txs {
		GetTx(v)
	}
}
