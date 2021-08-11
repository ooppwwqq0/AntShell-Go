package menu

import (
	"AntShell-Go/config"
	"AntShell-Go/models"
	"AntShell-Go/utils"
	"fmt"
	"os"
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
	limit = utils.IF(limit != 0, limit, DefaultLimit).(int)
	offset = utils.IF(offset != 0, offset, DefaultSize).(int)

	modeMax := menu.SttySize / count
	hostLen := len(hosts)
	pageMax := hostLen / offset

	// 处理菜单有几列
	var mode int
	if singleMode {
		mode = SingleMode
	} else {
		if customModeNum == 0 {
			mode = utils.IF(modeMax < DefaultModeMax, modeMax, DefaultMode).(int)
		} else {
			mode = utils.IF(customModeNum < modeMax, customModeNum, modeMax).(int)
		}
		mode = utils.IF(hostLen < mode, hostLen, mode).(int)
	}

	limit = utils.IF((limit < pageMax) && (limit > 0), limit, pageMax).(int)
	start := (limit-1)*offset + 1
	stop := limit*offset + 1

	lineFormat := "%5s %-24s %18s@%15s:%-5s "
	tailFormat := " All Pages %-5s %21s[c/C Clear] [n/N Back] Pages %-5s"
	headMsg := fmt.Sprintf(lineFormat, "[ID]", "NAME", "USER", "IP", "PORT")
	tailMsg := fmt.Sprintf(tailFormat, fmt.Sprintf("[%d]", hostLen), "", fmt.Sprintf("[%d]", 1))

	if !singleMode {
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
	stop = utils.IF(stop <= hostLen, stop, hostLen).(int)
	stop = utils.IF(singleMode, stop, hostLen).(int)

	// 输出主机条目
	for index, host := range hosts[start-1 : stop] {
		h := PrintInfo(index+1, host)
		if singleMode {
			fmt.Println(h)
		} else {
			rem := index % mode
			if mode == 1 || index+1 == hostLen {
				fmt.Println(h)
			} else if rem == 1 {
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
	hostPtr := models.NewHostPtr()

	switch argv.(type) {
	case int:
		num = utils.IF(argv.(int) != 0, argv.(int), num).(int)
	case string:
		hostPtr.SetSearch(argv.(string))
	}

	var hosts []models.Hosts
	hosts = hostPtr.Search(search, false, false)
	num = utils.IF(len(hosts) == 1, 1, num).(int)
	num = utils.IF(num <= len(hosts), num, 0).(int)
	if num != 0 {
		BannerPrint(menu.c)
		menu.Print(hosts, mode, DefaultLimit, DefaultSize, false)
		limit, offset := 0, utils.IF(customPage != 0, customPage, DefaultSize).(int)
		pageMax := len(hosts) / offset

		var input interface{}
		for num == 0 {
			fmt.Printf(
				ColorMsg("", BLUEL, false, false, BlankEnd),
				"\nInput your choose or [ 'q' | ctrl-c ] to quit!\nServer [ ID | IP | NAME ] >> ",
			)
			fmt.Scanln(&input)
			switch input.(type) {
			case string:
				switch strings.ToLower(input.(string)) {
				case "q", "quit", "exit":
					os.Exit(0)
				case "c", "clear":
					hostPtr.ClearSearch()
					continue
				case "n":
					limit = utils.IF(limit > 1, limit-1, 1).(int)
					continue
				case "m":
					limit = limit + 1
					continue
				default:
					if len(input.(string)) != 0 {
						hostPtr.SetSearch(input.(string))
						hosts = hostPtr.Search(input.(string), false, false)
					} else {
						limit = limit + 1
					}
				}
			case int:
				num = utils.IF(input.(int) <= len(hosts), input.(int), 0).(int)
				num = utils.IF(len(hosts) == 1, 1, num).(int)
			}
			limit = utils.IF(limit <= pageMax, limit, pageMax).(int)
		}
	}
	host = hosts[num-1]
	return
}
