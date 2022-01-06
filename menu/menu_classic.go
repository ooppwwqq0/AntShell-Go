package menu

import (
	"AntShell-Go/config"
	"AntShell-Go/models"
	"AntShell-Go/utils"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	DefaultSize    = 15
	DefaultLimit   = 1
	DefaultMode    = 5
	DefaultModeMax = 6
	SingleMode     = 1
	BlankEnd       = ""
	NewLineEnd     = "\n"
	count          = 76
)

type Classic struct {
	BaseMenu
}

func (menu *Classic) Init(c config.Config) {
	menu.c = c
	_, menu.SttySize = utils.GetSttySize()
}

func (menu *Classic) Print(hosts []models.Hosts, customModeNum int, limit int, offset int, singleMode bool) {
	offset = utils.IF(offset != 0, offset, DefaultSize).(int)
	limit = utils.IF(limit != 0, limit, DefaultLimit).(int)

	/*
		modeMax：通过终端宽度计算最大列数；
		hostLen：计算主机列表长度
		pageMax：通过主机列表长度计算最大页数
	*/
	modeMax := menu.SttySize / count
	hostLen := len(hosts)
	pageMax := int(math.Ceil(float64(hostLen) / float64(offset)))

	// 处理菜单列数
	var mode int
	if singleMode {
		mode = SingleMode
	} else {
		if customModeNum != 0 {
			// 如果用户指定了列数，并且没有超出最大列数，则展示用户指定列数
			mode = utils.IF(customModeNum < modeMax, customModeNum, modeMax).(int)
		}
		mode = utils.IF(mode <= 0, modeMax, mode).(int)
		// 如果列数超出了最大限制，则只展示最大列数
		mode = utils.IF(mode < DefaultModeMax, mode, DefaultModeMax).(int)
		// 如果主机列表长度小于列数，列数值为主机列表长度
		mode = utils.IF(hostLen < mode, hostLen, mode).(int)
		// 如果主机列表数据小于单页数量，则变为单列模式
		mode = utils.IF(hostLen < offset, SingleMode, mode).(int)
	}

	// 根据limit和offset处理主机记录开始，结束位置
	limit = utils.IF((limit <= pageMax) && (limit > 0), limit, pageMax).(int)
	start := (limit-1)*offset + 1
	stop := limit * offset

	stop = utils.IF(stop <= hostLen, stop, hostLen).(int)
	// 处理当mode=1时 输出模式和单列模式一致逻辑
	stop = utils.IF(singleMode || mode == SingleMode, stop, hostLen).(int)

	lineFormat := "%5s %-24s %18s@%15s:%-5s "
	tailFormat := " All Pages %-5s %21s[c/C Clear] [n/N Back] Pages %-5s"
	headMsg := fmt.Sprintf(lineFormat, "[ID]", "NAME", "USER", "IP", "PORT")
	tailMsg := fmt.Sprintf(tailFormat, fmt.Sprintf("[%d]", hostLen), "", fmt.Sprintf("[%d]", limit))

	if singleMode {
		BannerPrint(menu.c)
	}
	// 输出头部菜单
	var end string
	for i := 1; i <= mode; i++ {
		end = utils.IF(i == mode, NewLineEnd, BlankEnd).(string)
		ColorMsg(headMsg, WHITE, true, true, end)
		if i < mode {
			fmt.Print("  |  ", end)
		}
	}

	// 输出主机条目
	//for index, host := range hosts[start-1 : stop] {
	start = utils.IF(start > 1, start, 1).(int)
	for index := start - 1; index < stop; index++ {
		host := hosts[index]
		h := PrintInfo(index+1, host)
		if singleMode {
			fmt.Println(h)
		} else {
			rem := index % mode
			if mode == SingleMode || index+1 == hostLen {
				fmt.Println(h)
			} else if rem == mode-1 {
				fmt.Println(h)
			} else if rem < mode {
				fmt.Print(h, BlankEnd)
				fmt.Print("  |  ", BlankEnd)
			}
		}
	}

	// 输出尾部菜单
	for i := 1; i <= mode; i++ {
		end = utils.IF(i == mode, NewLineEnd, BlankEnd).(string)
		ColorMsg(utils.IF(singleMode, tailMsg, headMsg).(string), BLUEL, true, true, end)
		if i < mode {
			fmt.Print("  |  ", end)
		}
	}
}

func (menu *Classic) View(argv interface{}, num int, search string, mode int, customPage int) (host models.Hosts) {
	/*
		交互模式
			用户交互式选择主机信息
			当已知入参无法选取唯一主机记录时，进入交互模式
			交互模式下支持：主机id选择，主机名模糊匹配，主机ip模糊匹配，主机ip精确匹配；
			多次模糊匹配输入支持进一步筛选，支持清楚筛选记录
			交互模式下会自动进入单列模式，支持翻页
		无变量参数
			支持一位无变量参数，无变量参数可以直接进行精确搜索、模糊搜索、指定id选定主机记录
			通过无变量参数可以快速筛选主机记录，跳过交互模式，无变量参数可以和有变量参数配合使用
	*/
	hostPtr := models.NewHostPtr()

	// 无变量参数
	if argv != nil {
		if n, err := strconv.Atoi(argv.(string)); err == nil {
			num = utils.IF(n != 0, n, num).(int)
		} else {
			hostPtr.SetSearch(argv.(string))
		}
	}

	// 根据主机记录以及无变量参数尝试快速找到主机记录，或筛选主机记录
	var hosts []models.Hosts
	hosts = hostPtr.Search(search, false, false)
	num = utils.IF(len(hosts) == 1, 1, num).(int)
	num = utils.IF(num <= len(hosts), num, 0).(int)

	if num == 0 {
		BannerPrint(menu.c)
		menu.Print(hosts, mode, DefaultLimit, DefaultSize, false)
		limit, offset := 0, utils.IF(customPage != 0, customPage, DefaultSize).(int)
		pageMax := int(math.Ceil(float64(len(hosts)) / float64(offset)))

		var input string
		for num == 0 {
			fmt.Printf(
				ColorMsg("", BLUEL, false, false, BlankEnd),
				"\nInput your choose or [ 'q' | ctrl-c ] to quit!\nServer [ ID | IP | NAME ] >> ",
			)
			fmt.Scanln(&input)

			// 如果input转int时，err不为nil，则input是字符串
			if _, err := strconv.Atoi(input); err != nil {
				switch strings.ToLower(input) {
				case "q", "quit", "exit":
					os.Exit(0)
				case "c", "clear":
					hostPtr.ClearSearch()
					limit = 1
					hosts = hostPtr.Search("", false, false)
				case "n":
					limit = utils.IF(limit > 1, limit-1, 1).(int)
				case "m":
					limit = limit + 1
				default:
					if len(input) != 0 {
						hostPtr.SetSearch(input)
						hosts = hostPtr.Search(input, false, false)
					} else {
						limit = limit + 1
					}
				}
				input = ""
			} else {
				intInput, _ := strconv.Atoi(input)
				num = utils.IF(intInput <= len(hosts), intInput, 0).(int)
			}
			limit = utils.IF((limit <= pageMax) && (limit > 0), limit, pageMax).(int)
			//fmt.Println(hosts, mode, limit, offset)
			num = utils.IF(len(hosts) == 1, 1, num).(int)
			if num == 0 {
				menu.Print(hosts, mode, limit, offset, true)
			}
		}
	}
	hosts = hosts[num-1 : num]
	BannerPrint(menu.c)
	menu.Print(hosts, mode, DefaultLimit, DefaultSize, false)
	if len(hosts) >= 1 {
		host = hosts[0]
	}
	return
}
