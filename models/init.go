package models

import (
	"AntShell-Go/config"
	"AntShell-Go/utils"
	"github.com/astaxie/beego/orm"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	c := config.LoadConfig()

	orm.RegisterDriver("sqlite", orm.DRSqlite)
	orm.RegisterDataBase("default", "sqlite3", utils.ExpendUser(c.Default.DB_Path))
	orm.RunSyncdb("default", false, false)
}
