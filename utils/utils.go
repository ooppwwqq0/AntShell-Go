package utils

import (
	"os"
	"os/user"
	"path"
)

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

func ExpendUser(oldPath string) (realPath string) {
	if oldPath == "" {
		return
	}
	if oldPath[0:1] == "~" {
		user, _ := user.Current()
		realPath = path.Join(user.HomeDir, oldPath[1:])
	} else {
		realPath = oldPath
	}
	return
}
