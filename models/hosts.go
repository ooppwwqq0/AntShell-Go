package models

import (
	"AntShell-Go/utils"
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
	Path     string
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

// 获取全部数据
func (h *HostsPtr) GetAll() (rows []Hosts) {
	h.query.OrderBy("-sort").All(&h.Rows)
	rows = h.Rows
	return
}

// 搜索主机记录
func (h *HostsPtr) Search(search string, blankReturn bool, match bool) (rows []Hosts) {
	h.SetSearch(search)
	// 判断搜索内容是否是精确ip，如果是精确ip，会使用精确匹配模式
	match = utils.IF(utils.IsIP(search, true), true, match).(bool)
	if len(h.search) != 0 {
		cond := orm.NewCondition()
		if match {
			// 精确匹配ip记录
			for _, query := range h.search {
				cond = cond.Or("ip__in", query)
			}
		} else {
			// 模糊匹配ip、name记录，多个搜索条件and
			condOr := orm.NewCondition()
			for _, query := range h.search {
				cond = cond.AndCond(
					condOr.Or("ip__icontains", query).Or("name__icontains", query),
				)
			}
		}
		h.query.SetCond(cond).All(&h.Rows)
		rows = h.Rows
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

// 添加主机记录
func (h *HostsPtr) AddHost(host Hosts) (newHost Hosts) {
	cond := orm.NewCondition()
	cond = cond.And("ip", host.Ip).And("user", host.User).And("port", host.Port)
	h.query.SetCond(cond).All(&newHost)
	if newHost.Id != 0 {
		return
	}
	id, err := h.orm.Insert(&host)
	if err != nil {
		fmt.Println(err)
	}
	host.Id = int(id)
	return host
}

// 删除主机记录
func (h *HostsPtr) DelHost(host Hosts) (ok bool) {
	if _, err := h.orm.Delete(&host); err == nil {
		ok = true
	}
	return
}

// 修改主机记录
func (h *HostsPtr) UpdateHost(host Hosts) (n int64, err error) {
	n, err = h.orm.Update(&host)
	if err != nil {
		fmt.Println(err)
	}
	return
}

// 新增搜索词
func (h *HostsPtr) SetSearch(search string) {
	if search != "" && !utils.IsInArray(search, h.search) {
		h.search = append(h.search, search)
	}
}

// 清除搜索词
func (h *HostsPtr) ClearSearch() {
	h.search = []string{}
}

// 主机记录热度加1
func (h *HostsPtr) Sort(host Hosts, offset int) {
	host.Sort += offset
	h.orm.Update(&host, "sort")
}
