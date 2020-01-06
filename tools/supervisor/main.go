package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	hubble "github.com/vntchain/go-vnt"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/core/types"
	"github.com/vntchain/go-vnt/core/vm/supervisor"
	"github.com/vntchain/go-vnt/vntclient"
	"github.com/vntchain/vnt-explorer/models"
	"log"
	"math/big"
)

func main() {
	client, err := vntclient.Dial("ws://127.0.0.1:8546")
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress(supervisor.ContractAddr)
	query := hubble.FilterQuery{
		Addresses: []common.Address{contractAddress},
		FromBlock: big.NewInt(1),
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	abi, _ := supervisor.GetSuervisorABI()

	isMethod := func(name string, topic common.Hash) bool {
		if bytes.Equal(topic.Bytes(), common.BytesToHash(abi.Methods[name].Id()).Bytes()) {
			return true
		}
		return false
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			var (
				topic = vLog.Topics[0]
			)
			fmt.Println(topic)
			switch {
			case isMethod(supervisor.RegBizMeta, topic):
				d := &supervisor.BizMeta{}
				json.Unmarshal(vLog.Data, d)
				meta := models.BizMeta{
					No:        d.No,
					BizName:   d.BizName,
					BizType:   d.BizType,
					Desc:      d.Desc,
					Datas:     getJson(d.Datas),
					Tasks:     getJson(d.Tasks),
					Timestamp: d.Timestamp,
				}
				meta.Insert()
			case isMethod(supervisor.UpdateConfig, topic):
				d := &supervisor.Config{}
				json.Unmarshal(vLog.Data, d)
				fmt.Print(getJson(d))
			case isMethod(supervisor.RegisterBizContract, topic):
				d := &supervisor.BizContract{}
				json.Unmarshal(vLog.Data, d)
				c := models.BizContract{
					Address:   d.Address.Hex(),
					Owner:     d.Owner.Hex(),
					Name:      d.Name,
					Desc:      d.Desc,
					Status:    d.Status,
					TimeStamp: d.TimeStamp.Uint64(),
					BizNo:     d.BizNo,
				}
				c.Insert()
			case isMethod(supervisor.Report, topic):
				d := &supervisor.StructReport{}
				json.Unmarshal(vLog.Data, d)
				s := models.Report{
					ContractAddr: d.Addr,
					MetaNo:       uint64(d.MetaNo),
					Data:         getJson(d.Datas),
					BlockNumber:  vLog.BlockNumber,
					TxHash:       vLog.TxHash.Hex(),
					TimeStamp:    d.Timestamp,
				}
				s.Insert()
			}
		}
	}
}

func getJson(i interface{}) string {
	if bs, err := json.Marshal(i); err == nil {
		return string(bs)
	}
	return ""
}
