package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	registerModel()

	// beego.AppConfig.Set("mysql::user", "root")
	// beego.AppConfig.Set("mysql::pass", "root")
	// beego.AppConfig.Set("mysql::host", "127.0.0.1")
	// beego.AppConfig.Set("mysql::port", "3306")
	// beego.AppConfig.Set("mysql::db", "vnt")
	// beego.AppConfig.Set("node::rpc_host", "127.0.0.1")
	// beego.AppConfig.Set("node::rpc_port", "8545")

	var (
		dbuser = beego.AppConfig.String("mysql::user")
		dbpass = beego.AppConfig.String("mysql::pass")
		dbhost = beego.AppConfig.String("mysql::host")
		dbport = beego.AppConfig.String("mysql::port")
		dbname = beego.AppConfig.String("mysql::db")
	)

	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", dbuser, dbpass, dbhost, dbport, dbname)
	beego.Info("Will connect to mysql url", dbUrl)
	err := orm.RegisterDataBase("default", "mysql", dbUrl)
	if err != nil {
		beego.Error("failed to register database", err)
		panic(err.Error())
	}
}

func registerModel() {
	beego.Info("Will register models.")
	orm.RegisterModel(new(Account))
	orm.RegisterModel(new(Block))
	orm.RegisterModel(new(Node))
	orm.RegisterModel(new(TokenBalance))
	orm.RegisterModel(new(Transaction))
	orm.RegisterModel(new(Hydrant))
	orm.RegisterModel(new(MarketInfo))
	orm.RegisterModel(new(Subscription))
	orm.RegisterModel(new(Report))
	orm.RegisterModel(new(BizMeta))
	orm.RegisterModel(new(BizContract))
}
