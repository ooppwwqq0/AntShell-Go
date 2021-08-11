package models

import (
	"AntShell-Go/utils"
	"database/sql"
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
	Rows   []Hosts
	orm    orm.Ormer
	query  orm.QuerySeter
	search []string
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
func (h *HostsPtr) GetAll() (rows []Hosts) {
	h.query.All(&h.Rows)
	rows = h.Rows
	return
}

func (h *HostsPtr) Search(search string, blankReturn bool, match bool) (rows []Hosts) {
	h.SetSearch(search)
	match = utils.IF(utils.IsIP(search, true), true, match).(bool)
	if len(h.search) != 0 {
		if match {
			//rows = h.query.Filter()
			//hosts = SESSION.query(Hosts).filter(Hosts.ip.in_(list(map(lambda x: "%s" % x, self.search)))).all()

		} else {
			//hosts = SESSION.query(Hosts).filter(
			//	or_(*list(map(lambda x: Hosts.ip.like('%%{0}%%'.format(x)), self.search)),
			//*list(map(lambda x: Hosts.name.like('%%{0}%%'.format(x)), self.search)))).all()
		}
	} else {
		rows = h.GetAll()
	}
	if len(rows) == 0 && !blankReturn {
		rows = h.GetAll()
		if len(rows) == 0 {
			rows = []Hosts{}
		}
	}

	return
}

func (h *HostsPtr) SetSearch(search string) {
	if search != "" && !utils.IsInArray(search, h.search) {
		h.search = append(h.search, search)
	}
}

func (h *HostsPtr) ClearSearch() {
	h.search = []string{}
}
