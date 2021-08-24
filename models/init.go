package models

import (
	"AntShell-Go/config"
	"AntShell-Go/utils"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/mitchellh/go-homedir"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

func Init(c config.Config) {
	dbPath, _ := homedir.Expand(c.Default.DB_Path)
	if !utils.IsFile(dbPath) {
		logs.Info("创建默认数据文件路径:", path.Dir(dbPath))
		os.MkdirAll(path.Dir(dbPath), 0755)
	}

	orm.RegisterDriver("sqlite", orm.DRSqlite)
	orm.RegisterDataBase("default", "sqlite3", dbPath)
	orm.RunSyncdb("default", false, false)
}
