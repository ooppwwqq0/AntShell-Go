package utils

import (
	"github.com/astaxie/beego/logs"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// 判断文件是否存在
func IsExist(fileAddr string) bool {
	// 读取文件信息，判断文件是否存在
	_, err := os.Stat(fileAddr)
	if err != nil {
		if os.IsExist(err) { // 根据错误类型进行判断
			return true
		}
		return false
	}
	return true
}

// 判断目录是否存在
func IsDir(fileAddr string) bool {
	s, err := os.Stat(fileAddr)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断文件是否存在
func IsFile(fileAddr string) bool {
	s, err := os.Stat(fileAddr)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

// 实现简单的三元运算
func IF(condition bool, trueVal interface{}, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

// 清屏
func Clear() {
	optSys := runtime.GOOS
	if optSys == "linux" || optSys == "darwin" {
		//执行clear指令清除控制台
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		//执行clear指令清除控制台
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

// 获取终端大小
func GetSttySize() (high int, width int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, _ := cmd.Output()

	outList := strings.Split(strings.Replace(string(out), "\n", "", 1), " ")
	high, _ = strconv.Atoi(outList[0])
	width, _ = strconv.Atoi(outList[1])
	return high, width
}

// 堡垒机根据totp码获取动态码
func GetPasswdByTotp(totp string) (passwd string) {
	cmd := exec.Command("oathtool", "-b", "--totp", totp)
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		logs.Error(err)
		os.Exit(1)
	}
	passwd = strings.Replace(string(out), "\n", "", 1)
	return
}

func ExecCommand(c string) (output string, err error) {
	cmd := exec.Command(c)
	cmd.Stdin = os.Stdin
	var out []byte
	out, err = cmd.Output()
	if err != nil {
		logs.Error(err)
		os.Exit(1)
	}
	output = string(out)
	return
}

// 判断值是否在对象中
func IsInArray(item interface{}, array interface{}) (bolFind bool) {

	switch array.(type) {
	case []string:
		data := array.([]string)
		strItem, ok := item.(string)
		if !ok {
			bolFind = false
			return bolFind
		}
		for _, one := range data {
			if one == strItem {
				bolFind = true
				return bolFind
			}
		}

	case []int:
		data := array.([]int)
		intItem, ok := item.(int)
		if !ok {
			bolFind = false
			return bolFind
		}
		for _, one := range data {
			if one == intItem {
				bolFind = true
				return bolFind
			}
		}
	}

	bolFind = false
	return bolFind
}

// 粗糙实现判断是否ip
func IsIP(ipAddr string, mask bool) bool {
	if ipAddr == "" {
		return false
	}
	ipList := strings.Split(ipAddr, ".")
	l := IF(mask, 4, len(ipList)).(int)
	l = IF(l <= 4, l, 4).(int)
	var qi []int
	for _, s := range ipList {
		ip, err := strconv.Atoi(s)
		if err != nil {
			return false
		} else if ip >= 0 && ip <= 255 {
			qi = append(qi, ip)
		}
	}
	return len(qi) == l
}

func ReadByFile(path string) (context string, err error) {
	realPath, _ := homedir.Expand(path)
	if IsFile(realPath) {
		var f []byte
		f, err = ioutil.ReadFile(realPath)
		if err != nil {
			logs.Error(err)
		}
		context = string(f)
	}
	return
}

func CreateAndWrite(path string, context string) (err error) {
	realPath, _ := homedir.Expand(path)
	if IsFile(realPath) {
		err = os.Remove(realPath)
		if err != nil {
			logs.Error(err)
		}
	}
	f, err := os.Create(realPath)
	if err != nil {
		logs.Error(err)
	}
	f.WriteString(context)
	f.Close()

	return
}
