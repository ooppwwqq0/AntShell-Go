package config

import (
	"AntShell-Go/utils"
	"github.com/ooppwwqq0/goconfig"
	"os"
	"path"
)

const (
	INI_TYPE     = "ini"
	YAML_TYPE    = "yaml"
	DEFAULT_PATH = "~/.antshell"
	CONFIG_NAME  = "antshell.cfg"
	ENV_PATH     = "ANTSHELL_CONFIG"
)

type Config struct {
	Default DefaultSection `ini:"default"`
	User    UserSection    `ini:"user"`
	Bastion BastionSection `ini:"BASTION"`
}

type DefaultSection struct {
	DB_Path      string `ini:"DB_PATH"`
	Key_Path     string `ini:"KEY_PATH"`
	LangSet      string `ini:"LANGSET"`
	Banner_Color string `ini:"BANNER_COLOR"`
	Debug        string `ini:"DEBUG"`
	Page         string `ini:"PAGE"`
	Engine       string `ini:"ENGINE"`
}

type UserSection struct {
	UserName string `ini:"USERNAME"`
	Password string `ini:"PASSWORD"`
	Port     string `ini:"PORT"`
}

type BastionSection struct {
	Bastion_Host          string `ini:"BASTION_HOST"`
	Bastion_Port          string `ini:"BASTION_PORT"`
	Bastion_User          string `ini:"BASTION_USER"`
	Bastion_Passwd_Prefix string `ini:"BASTION_PASSWD_PREFIX"`
	Bastion_Passwd        string `ini:"BASTION_PASSWD"`
	Bastion_Totp          string `ini:"BASTION_TOTP"`
}

func FindConfig() (configPath string) {

	var pathList []string
	customPath := os.Getenv(ENV_PATH)
	if utils.IsDir(customPath) {
		pathList = append(pathList, path.Join(customPath, CONFIG_NAME))
	}
	pathList = append(pathList, path.Join(DEFAULT_PATH, CONFIG_NAME))
	pathList = append(pathList, path.Join("/etc/antshell/", CONFIG_NAME))
	for _, config := range pathList {
		configPath = utils.ExpendUser(config)
		if utils.IsFile(configPath) {
			break
		}
	}
	return
}

func LoadConfig() (config Config) {
	configPath := FindConfig()
	cfg, err := goconfig.LoadConfigFile(configPath)
	if err != nil {
		panic("错误")
	}
	err = cfg.Decode(&config)
	return
}
