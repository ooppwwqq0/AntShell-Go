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
	"reflect"
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
	flag.BoolVar(&option.Host.Delete, "d", false, lang["delete"])
	flag.StringVar(&option.Host.Name, "name", "", lang["name"])
	flag.StringVar(&option.Host.User, "user", "", lang["user"])
	flag.StringVar(&option.Host.Passwd, "passwd", "", lang["passwd"])
	flag.IntVar(&option.Host.Port, "port", 0, lang["port"])
	flag.StringVar(&option.Host.Sudo, "sudo", "", lang["sudo"])
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
	host = models.Hosts{
		Ip:      option.Host.Add,
		Sudo:    option.Host.Sudo,
		Name:    utils.IF(option.Host.Name != "", option.Host.Name, option.Host.Add).(string),
		User:    utils.IF(option.Host.User != "", option.Host.User, c.User.UserName).(string),
		Passwd:  utils.IF(option.Host.Passwd != "", option.Host.Passwd, c.User.Password).(string),
		Port:    utils.IF(option.Host.Port != 0, option.Host.Port, c.User.Port).(int),
		Bastion: utils.IF(option.Manager.Bastion, engine.BastionOn, engine.BastionOff).(int),
	}
	return
}

func GetUserInput(title string, value interface{}, defaultValue interface{}, valueType string, isNone bool) (newValue interface{}) {
	/*
		编辑主机信息用户交互
		title: 显示标题
		value: 命令行参数值，优先于default值
		default: 默认值，原数据库中数据
		vType: 指定用户输入值
		isNone: 用户输入是否可为空，默认False，用户输入为空时使用默认值
	*/
	msg := menu.ColorMsg("", menu.BLUEL, false, false, "")

	switch valueType {
	case "string":
		var input string
		fmt.Printf(msg, fmt.Sprintf("New %s [defalut: %s] [new: %s] >> ", title, defaultValue, value))
		fmt.Scanln(&input)
		if input == "" && !isNone {
			newValue = utils.IF(value != "", value, defaultValue)
		} else {
			newValue = input
		}
	case "int":
		var input int
		value = utils.IF(value != 0 && !isNone, strconv.Itoa(value.(int)), "")
		fmt.Printf(msg, fmt.Sprintf("New %s [defalut: %d] [new: %s] >> ", title, defaultValue, value))
		fmt.Scanln(&input)
		if input == 0 && !isNone {
			newValue = utils.IF(value != "", value, defaultValue)
			fmt.Println("ddd")
		} else {
			newValue = input
		}
		newValue = strconv.Itoa(newValue.(int))

	}
	fmt.Println(newValue)
	fmt.Println(reflect.TypeOf(newValue))

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
		fmt.Println(newHost)
		menu.ColorMsg("Confirm [ y|n ] >> ", menu.GREEN, true, false, "")
		fmt.Scanln(&input)

		if strings.ToLower(input) == "y" {
			fmt.Println(newHost)
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
	hostPtr := models.NewHostPtr()
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
		fmt.Println(host)
		newHost := GetHostByUser(host)
		fmt.Println(newHost)
		os.Exit(0)
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
	client.Connection(option.Host.Sudo)
	os.Exit(0)
}
