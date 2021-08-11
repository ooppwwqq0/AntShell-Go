package main

import (
	"AntShell-Go/config"
	"AntShell-Go/engine"
	"AntShell-Go/menu"
	"AntShell-Go/models"
	"AntShell-Go/utils"
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Option struct {
	Host    HostOption
	Manager ManagerOption
}

type HostOption struct {
	Add    string
	Edit   bool
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
	Bastion bool
	Version bool
	Argv    interface{}
}

var (
	c      = config.LoadConfig()
	option Option
	client engine.ClientSSH
)

func init() {
	lang := utils.LANG[c.Default.LangSet]
	flag.BoolVar(&option.Manager.List, "l", false, lang["list"])
	flag.IntVar(&option.Manager.Mode, "m", 0, lang["mode"])
	flag.IntVar(&option.Manager.Num, "n", 0, lang["num"])
	flag.StringVar(&option.Manager.Search, "s", "", lang["search"])
	flag.BoolVar(&option.Manager.Bastion, "B", false, lang["bastion"])
	flag.BoolVar(&option.Manager.Version, "version", false, lang["version"])

	flag.StringVar(&option.Host.Add, "a", "", lang["add"])
	flag.BoolVar(&option.Host.Edit, "e", false, lang["edit"])
	flag.StringVar(&option.Host.Delete, "d", "", lang["delete"])
	flag.StringVar(&option.Host.Name, "name", "", lang["name"])
	flag.StringVar(&option.Host.User, "user", "", lang["user"])
	flag.StringVar(&option.Host.Passwd, "passwd", "", lang["passwd"])
	flag.IntVar(&option.Host.Port, "port", 0, lang["port"])
	flag.StringVar(&option.Host.Sudo, "sudo", "", lang["sudo"])
	flag.Usage = usage
	flag.Parse()
	if len(flag.Args()) != 0 {
		option.Manager.Argv = flag.Args()[0]
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `AntShell version: AntShell/1.0
Usage: antshell|a [ -h | --version ] [-l [-m 2] ] [ v | -n 1 | -s 'ip|name' ] [ -A ] [ -B ]
        [ -e | -d ip | -a ip [--name tag | --user root | --passwd *** | --port 22 | --sudo root ] ]
`)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, `
# Add host record	
	a -a 10.0.0.1
	a -a 10.0.0.1 -name app01
	a -a 10.0.0.1 -name app01 -passwd 123456
	a -a 10.0.0.1 -name app01 -user root -passwd 123456 
	a -a 10.0.0.1 -name app01 -user root -passwd 123456 -port 22022
	a -a 10.0.0.1 -name app01 -user ubuntu -passwd 123456 -sudo root
	a -a 10.0.0.1 -name app01 -user ubuntu -passwd 123456 -port 22022 -sudo root -B
# Delete host record
	a -d 10.0.0.1
	a -d app01
# Edit host record
	a -e
	a -e -s 10.0.0.1
	a -e -s app01 -n 2
# List host record
	a -l
	a -l -m 2
# Login host
	a
	a 2
	a app01
	a 10.0.0.0.1
	a app01 -n 2
	a -s 10.0.0.1 -n 1
	a -s app01 -n 2
`)
}

func main() {
	if option.Manager.Version {
		fmt.Printf("%s %s\n", utils.ProgramName, utils.Version)
		os.Exit(0)
	}
	hostPtr := models.NewHostPtr()
	m := menu.New(c)

	switch {
	case option.Manager.List:
		hosts := hostPtr.GetAll()
		m.Print(hosts, option.Manager.Mode, 1, 15, false)
		os.Exit(0)
	case option.Host.Add != "":
	case option.Host.Edit:
		break
	case option.Host.Delete != "":
		break
	default:
		//m.Print(hosts, option.Manager.Mode, 1, 15, false)
		customPage, _ := strconv.Atoi(c.Default.Page)
		host := m.View(
			option.Manager.Argv,
			option.Manager.Num, option.Manager.Search,
			option.Manager.Mode, customPage,
		)
		menu.BannerPrint(c)
		m.Print([]models.Hosts{host}, 1, 1, 1, false)

		client.Init(host, c)
		client.Connection(option.Host.Sudo)
		os.Exit(0)
	}
}
