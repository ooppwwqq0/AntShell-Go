package models

import (
	"AntShell-Go/config"
	"github.com/astaxie/beego/orm"
	"github.com/mitchellh/go-homedir"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	c := config.LoadConfig()

	orm.RegisterDriver("sqlite", orm.DRSqlite)
	dbPath, _ := homedir.Expand(c.Default.DB_Path)
	orm.RegisterDataBase("default", "sqlite3", dbPath)
	orm.RunSyncdb("default", false, false)
}
