package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/vntchain/vnt-explorer/models"
)

func main() {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RunSyncdb("default", true, true)
	alterTable()
}

func registerModel() {
	beego.Info("Will register models.")
	orm.RegisterModel(new(models.Account))
	orm.RegisterModel(new(models.Block))
	orm.RegisterModel(new(models.Node))
	orm.RegisterModel(new(models.TokenBalance))
	orm.RegisterModel(new(models.Transaction))
}

func alterTable() {
	needAlterMap := make(map[string][]string)
	needAlterMap["account"] = []string{"balance", "token_amount", "token_acct_count"}
	needAlterMap["block"] = []string{"number"}
	needAlterMap["node"] = []string{"votes"}
	needAlterMap["token_balance"] = []string{"balance"}
	for tableName, columns := range needAlterMap {
		for _, col := range columns {
			var err error
			if tableName == "block" && col == "number" {
				err = alterColumn(tableName, col, "decimal(64,0) PRIMARY KEY")
			} else {
				err = alterColumn(tableName, col, "decimal(64,0)")
			}
			if err != nil {
				fmt.Println(err)
			}
		}
	}

}

func alterColumn(tableName, column, dataType string) error {
	// TODO 直接ALTER COLUMN会报错，待查证
	o := orm.NewOrm()
	alterString := fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", tableName, column)
	_, err := o.Raw(alterString).Exec()
	if err != nil {
		return fmt.Errorf("ALTER TABLE %s DROP COLUMN error: ", tableName, err)
	}
	alterString = fmt.Sprintf("ALTER TABLE %s ADD %s %s", tableName, column, dataType)
	_, err = o.Raw(alterString).Exec()
	if err != nil {
		return fmt.Errorf("ALTER TABLE %s ADD  error: ", tableName, err)
	}
	return nil
}
