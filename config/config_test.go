package config

import (
	"fmt"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	conf, _ := LoadConfig()
	for k := range conf.Bastion_List {
		fmt.Println(k, conf.Bastion_List[k])
	}
	fmt.Println(conf)
}
