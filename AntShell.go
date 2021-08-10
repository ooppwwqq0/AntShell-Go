package main

import (
	"AntShell-Go/config"
	"AntShell-Go/models"
	"AntShell-Go/utils"
	"flag"
	"fmt"
	"os"
)

type Option struct {
	Host    HostOption
	Manager ManagerOption
}

type HostOption struct {
	Add    string
	Edit   string
	Delete string
	Name   string
	User   string
	Passwd string
	Port   int
	Sudo   string
}

type ManagerOption struct {
	List    bool
	Mode    int
	Num     int
	Search  string
	Agent   bool
	Bastion bool
	Engine  string
}

var (
	option Option
)

func init() {
	c := config.LoadConfig()
	lang := utils.LANG[c.Default.LangSet]
	flag.StringVar(&option.Host.Add, "-a", "", lang["add"])
	flag.Usage = usage
	flag.Parse()
}

func usage() {
	fmt.Fprintf(os.Stderr, `准入查询工具 version: get_info_by_cloudapp/1.0
Usage: get_info_by_cloudapp [-h] [-bns app.product] [-app app] [-sub sub] [-product product] [-path product_sub] [-level 1]

`)
	flag.PrintDefaults()
}

func main() {
	//engine.SSH()
	//c := config.LoadConfig()
	//fmt.Println(c)
	//models.Sql()
	host := models.NewHostPtr()
	host.GetAll()

}
