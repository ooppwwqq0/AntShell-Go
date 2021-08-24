package main

import (
	"AntShell-Go/config"
	"AntShell-Go/engine"
	"AntShell-Go/menu"
	"AntShell-Go/models"
	"AntShell-Go/utils"
	"flag"
	"fmt"
	"github.com/astaxie/beego/logs"
	"os"
	"strconv"
	"strings"
)

type Option struct {
	Host    HostOption
	Manager ManagerOption
}

type HostOption struct {
	Add    string
	Edit   bool
	Delete bool
	Name   string
	User   string
	Passwd string
	Port   int
	Sudo   string
	Path   string
}

type ManagerOption struct {
	List    bool
	Mode    int
	Num     int
	Search  string
	Bastion bool
	Version bool
	Argv    interface{}
	Totp    bool
}

var (
	c      config.Config
	option Option
	client engine.ClientSSH
)

func init() {
	var err error
	c, err = config.LoadConfig()
	if err != nil {
		config.InitConfig()
		c, err = config.LoadConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	models.Init(c)

	lang := utils.LANG[c.Default.LangSet]
	flag.BoolVar(&option.Manager.List, "l", false, lang["list"])
	flag.IntVar(&option.Manager.Mode, "m", 0, lang["mode"])
	flag.IntVar(&option.Manager.Num, "n", 0, lang["num"])
	flag.StringVar(&option.Manager.Search, "s", "", lang["search"])
	flag.BoolVar(&option.Manager.Bastion, "B", false, lang["bastion"])
	flag.BoolVar(&option.Manager.Version, "version", false, lang["version"])
	flag.BoolVar(&option.Manager.Totp, "totp", false, lang["totp"])

	flag.StringVar(&option.Host.Add, "a", "", lang["add"])
	flag.BoolVar(&option.Host.Edit, "e", false, lang["edit"])
	flag.BoolVar(&option.Host.Delete, "d", false, lang["delete"])
	flag.StringVar(&option.Host.Name, "name", "", lang["name"])
	flag.StringVar(&option.Host.User, "user", "", lang["user"])
	flag.StringVar(&option.Host.Passwd, "passwd", "", lang["passwd"])
	flag.IntVar(&option.Host.Port, "port", 0, lang["port"])
	flag.StringVar(&option.Host.Sudo, "sudo", "", lang["sudo"])
	flag.StringVar(&option.Host.Path, "path", "", lang["path"])
	flag.Usage = usage
	flag.Parse()
	if len(flag.Args()) != 0 {
		option.Manager.Argv = flag.Args()[0]
		/*
			当flag遇到non-flag时会停止继续解析，将从non-flag开始的所有参数认定为non-flag
			这时所有后面的参数都不能正常运行
			通过源码分析可以发现flag.Parse实际是执行了flag.CommandLine.Parse方法
			那我们就可以通过以下方法让其他参数继续解析
		*/
		flag.CommandLine.Parse(flag.Args()[1:])
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

func GetHostByConfig(c config.Config) (host models.Hosts) {
	defaultPort, _ := strconv.Atoi(c.User.Port)
	host = models.Hosts{
		Ip:      option.Host.Add,
		Sudo:    option.Host.Sudo,
		Name:    utils.IF(option.Host.Name != "", option.Host.Name, option.Host.Add).(string),
		User:    utils.IF(option.Host.User != "", option.Host.User, c.User.UserName).(string),
		Passwd:  utils.IF(option.Host.Passwd != "", option.Host.Passwd, c.User.Password).(string),
		Port:    utils.IF(option.Host.Port != 0, option.Host.Port, defaultPort).(int),
		Bastion: utils.IF(option.Manager.Bastion, engine.BastionOn, engine.BastionOff).(int),
		Path:    utils.IF(option.Host.Path != "", option.Host.Path, c.User.Path).(string),
	}
	return
}

func GetUserInputStr(title string, value string, defaultValue string, isNone bool) (newValue string) {
	/*
		编辑主机信息用户交互
		title: 显示标题
		value: 命令行参数值，优先于default值
		default: 默认值，原数据库中数据
		vType: 指定用户输入值
		isNone: 用户输入是否可为空，默认False，用户输入为空时使用默认值
	*/
	msg := menu.ColorMsg("", menu.BLUEL, false, false, "")

	var input string
	fmt.Printf(msg, fmt.Sprintf("New %s [defalut: %s] [new: %s] >> ", title, defaultValue, value))
	fmt.Scanln(&input)
	if input == "" {
		newValue = utils.IF(value != "", value, defaultValue).(string)
	} else {
		newValue = input
	}

	return
}

func GetUserInputInt(title string, value int, defaultValue int, isNone bool) (newValue int) {
	/*
		编辑主机信息用户交互
		title: 显示标题
		value: 命令行参数值，优先于default值
		default: 默认值，原数据库中数据
		vType: 指定用户输入值
		isNone: 用户输入是否可为空，默认False，用户输入为空时使用默认值
	*/
	msg := menu.ColorMsg("", menu.BLUEL, false, false, "")

	var input int
	valueStr := utils.IF(value != 0 && !isNone, strconv.Itoa(value), "").(string)
	fmt.Printf(msg, fmt.Sprintf("New %s [defalut: %d] [new: %s] >> ", title, defaultValue, valueStr))
	fmt.Scanln(&input)
	if input == 0 {
		newValue = utils.IF(value != 0 && !isNone, value, defaultValue).(int)
	} else {
		newValue = input
	}
	return
}

func GetHostByUser(host models.Hosts) (newHost models.Hosts) {
	for {
		var input string
		newHost = host
		newHost.Name = GetUserInputStr("Name", option.Host.Name, host.Name, false)
		newHost.User = GetUserInputStr("User", option.Host.User, host.User, false)
		newHost.Passwd = GetUserInputStr("Passwd", option.Host.Passwd, host.Passwd, false)
		newHost.Port = GetUserInputInt("Port", option.Host.Port, host.Port, false)
		newHost.Sudo = GetUserInputStr("Sudo", option.Host.Sudo, host.Sudo, false)
		optionBastion := utils.IF(option.Manager.Bastion, engine.BastionOn, engine.BastionOff).(int)
		newHost.Bastion = GetUserInputInt("Bastion", optionBastion, host.Bastion, true)
		newHost.Path = GetUserInputStr("Path", option.Host.Path, host.Path, false)

		menu.ColorMsg("Confirm [ y|n ] >> ", menu.GREEN, true, false, "")
		fmt.Scanln(&input)

		if strings.ToLower(input) == "y" {
			break
		}

	}
	return
}

func main() {
	if option.Manager.Version {
		fmt.Printf("%s %s\n", utils.ProgramName, utils.Version)
		os.Exit(0)
	}

	if option.Manager.Totp {
		fmt.Println(utils.GetPasswdByTotp(c.Bastion.Bastion_Totp))
		os.Exit(0)
	}
	hostPtr := models.NewHostPtr()
	hosts := hostPtr.GetAll()
	if len(hosts) == 0 {
		logs.Warn("Please Add Host Record!")
		logs.Info("a -help")
		os.Exit(1)
	}

	m := menu.New(c)

	var host models.Hosts
	switch {
	case option.Manager.List:
		hosts := hostPtr.GetAll()
		menu.BannerPrint(c)
		m.Print(hosts, option.Manager.Mode, menu.DefaultLimit, menu.DefaultSize, false)
		os.Exit(0)
	case option.Host.Add != "":
		if !utils.IsIP(option.Host.Add, true) {
			menu.ColorMsg("wrong ip: "+option.Host.Add, menu.RED, true, false, "\n")
			os.Exit(1)
		}
		host = GetHostByConfig(c)
		host = hostPtr.AddHost(host)
	}
	if host.Id == 0 {
		customPage, _ := strconv.Atoi(c.Default.Page)
		host = m.View(
			option.Manager.Argv,
			option.Manager.Num, option.Manager.Search,
			option.Manager.Mode, customPage,
		)
	}

	switch {
	case option.Host.Edit:
		host = GetHostByUser(host)
		_, err := hostPtr.UpdateHost(host)
		if err != nil {
			os.Exit(1)
		}
	case option.Host.Delete:
		var input string
		for {
			msg := "Confirm [ y|n ] >> "
			fmt.Printf(menu.ColorMsg("", menu.GREEN, false, false, ""), msg)
			fmt.Scanln(&input)
			switch strings.ToLower(input) {
			case "y":
				hostPtr.DelHost(host)
				menu.ColorMsg(fmt.Sprintf("Delete IP: %s Success!", host.Ip), menu.BLUE, true, false, "\n")
				os.Exit(0)
			case "n":
				os.Exit(0)
			}
		}
	}

	client.Init(host, c)
	client.Connection(option.Host.Sudo, option.Host.Path)
	os.Exit(0)
}
