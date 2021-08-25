package menu

import (
	"AntShell-Go/config"
	"AntShell-Go/models"
	"AntShell-Go/utils"
	"fmt"
)

const (
	BLACK = iota + 30
	RED
	GREEN
	YELLOW
	BLUE
	PINK
	BLUEL
	WHITE
)

var (
	ColorCode = map[string]int{
		"red":    RED,
		"green":  GREEN,
		"yellow": YELLOW,
		"blue":   BLUE,
		"pink":   PINK,
		"cblue":  BLUEL,
		"white":  WHITE,
	}
	ColorOffset = 10
)

type BaseMenu struct {
	Menu
	SttySize int
	c        config.Config
}

type Menu interface {
	Init(c config.Config)
	Print(hosts []models.Hosts, customModeNum int, limit int, offset int, singleMode bool)
	View(argv interface{}, num int, search string, mode int, page int) (host models.Hosts)
}

func New(c config.Config) (menu Menu) {
	menu = &Classic{}
	menu.Init(c)
	return
}

func ColorMsg(msg string, color int, print bool, isTitle bool, end string) string {
	titleSign := BLACK
	if isTitle {
		titleSign = color + ColorOffset
		color = BLACK
	}
	colorSign := fmt.Sprintf("%c[1;%d;%dm", 0x1B, titleSign, color) + "%s" + fmt.Sprintf("%c[0m", 0x1B)
	if print {
		colorMsg := fmt.Sprintf(colorSign, msg)
		fmt.Print(colorMsg, end)
		return colorMsg
	}
	return colorSign
}

func PrintInfo(index int, host models.Hosts) (infoFormat string) {
	userName := utils.IF(host.Sudo != "", host.Sudo, host.User).(string)
	userColor := utils.IF(host.Sudo != "", YELLOW, GREEN).(int)
	bastionColor := utils.IF(host.Bastion == 1, YELLOW, GREEN).(int)
	infoFormat = fmt.Sprintf(
		"%s %s %s@%s:%s ",
		fmt.Sprintf(
			ColorMsg("", YELLOW, false, false, BlankEnd),
			fmt.Sprintf("%5s", fmt.Sprintf("[%d]", index)),
		),
		fmt.Sprintf(
			ColorMsg("", bastionColor, false, false, BlankEnd),
			fmt.Sprintf("%-24s", host.Name),
		),
		fmt.Sprintf(
			ColorMsg("", userColor, false, false, BlankEnd),
			fmt.Sprintf("%18s", userName),
		),
		fmt.Sprintf(
			ColorMsg("", GREEN, false, false, BlankEnd),
			fmt.Sprintf("%15s", host.Ip),
		),
		fmt.Sprintf(
			ColorMsg("", GREEN, false, false, BlankEnd),
			fmt.Sprintf("%-5d", host.Port),
		),
	)
	return
}

func BannerPrint(c config.Config) {
	utils.Clear()
	ColorMsg(utils.GetBanner(c.Default.Banner), ColorCode[c.Default.Banner_Color], true, false, NewLineEnd)
}
