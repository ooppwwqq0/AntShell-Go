package models

import (
	"database/sql"
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

var (
	db *sql.DB
)

type Hosts struct {
	Id       int
	Sort     int
	Name     string
	Ip       string
	User     string
	Passwd   string
	Port     int
	Sudo     string
	Bastion  int
	CreateAt time.Time
	UpdateAt time.Time
}

type HostsPtr struct {
	Rows  []Hosts
	orm   orm.Ormer
	query orm.QuerySeter
}

func NewHostPtr() *HostsPtr {
	host := &HostsPtr{}
	host.init()
	return host
}

func init() {
	orm.RegisterModel(new(Hosts))
}

func (h *HostsPtr) init() {
	h.orm = orm.NewOrm()
	h.orm.Using("default")
	h.query = h.orm.QueryTable(new(Hosts))
}
func (h *HostsPtr) GetAll() {
	h.query.All(&h.Rows)
	for _, line := range h.Rows {
		fmt.Println(line)
	}
}
