package config

import (
	"AntShell-Go/utils"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/mitchellh/go-homedir"
	"github.com/ooppwwqq0/goconfig"
	"os"
	"path"
)

const (
	TypeINI     = "ini"
	TypeYaml    = "yaml"
	DefaultPath = "~/.antshell"
	EtcPath     = "/etc/antshell/"
	EnvPath     = "ANTSHELL_CONFIG"
	ConfName    = "antshell.cfg"
)

type Config struct {
	Default      DefaultSection    `ini:"default"`
	User         UserSection       `ini:"user"`
	Bastion      BastionSection    `ini:"bastion"`
	Bastion_List map[string]string `ini:"bastion_list"`
}

type DefaultSection struct {
	DB_Path      string `ini:"DB_PATH"`
	Key_Path     string `ini:"KEY_PATH"`
	LangSet      string `ini:"LANGSET"`
	Banner_Color string `ini:"BANNER_COLOR"`
	Debug        string `ini:"DEBUG"`
	Page         string `ini:"PAGE"`
	Backup_Dir   string `ini:"BACKUP_DIR"`
	Banner       string `ini:"BANNER"`
}

type UserSection struct {
	UserName string `ini:"USERNAME"`
	Password string `ini:"PASSWORD"`
	Port     string `ini:"PORT"`
	Path     string `ini:"PATH"`
}

type BastionSection struct {
	Bastion_Host          string `ini:"BASTION_HOST"`
	Bastion_Port          string `ini:"BASTION_PORT"`
	Bastion_User          string `ini:"BASTION_USER"`
	Bastion_Passwd_Prefix string `ini:"BASTION_PASSWD_PREFIX"`
	Bastion_Passwd        string `ini:"BASTION_PASSWD"`
	Bastion_Totp          string `ini:"BASTION_TOTP"`
}

/*
获取配置文件路径
优先级：环境变量（ANTSHELL_CONFIG）> DefaultPath > EtcPath
*/
func FindConfig() (configPath string, err error) {

	var pathList []string
	customPath := os.Getenv(EnvPath)
	if utils.IsDir(customPath) {
		pathList = append(pathList, path.Join(customPath, ConfName))
	}
	pathList = append(pathList, path.Join(DefaultPath, ConfName))
	pathList = append(pathList, path.Join(EtcPath, ConfName))
	for _, config := range pathList {
		configPath, _ = homedir.Expand(config)
		if utils.IsFile(configPath) {
			return configPath, err
		}
	}
	err = errors.New("找不到配置文件")
	return
}

// 加载配置文件
func LoadConfig() (config Config, err error) {
	configPath, err := FindConfig()
	if err != nil {
		return config, err
	}

	var cfg *goconfig.ConfigFile
	cfg, err = goconfig.LoadConfigFile(configPath)
	if err != nil {
		return config, err
	}
	err = cfg.Decode(&config)
	return
}

// 初始化配置文件
func InitConfig() {
	logs.Info("开始初始化配置文件")
	defaultPath, _ := homedir.Expand(DefaultPath)
	if !utils.IsExist(defaultPath) {
		logs.Info("创建默认配置文件路径:", defaultPath)
		os.MkdirAll(defaultPath, 0755)
	}
	if !utils.IsFile(path.Join(defaultPath, ConfName)) {
		logs.Info("创建默认配置文件:", path.Join(defaultPath, ConfName))
		err := RestoreAssets(defaultPath, ConfName)
		logs.Error(err)
	}
}
