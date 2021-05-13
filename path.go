package golibs

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

// HomePath 获取当前用户根目录
func HomePath() (string, error) {
	user, err := user.Current()
	if nil == err {
		return user.HomeDir, nil
	}

	// cross compile support
	if "windows" == runtime.GOOS {
		return homeWindows()
	}

	// Unix-like system, so just assume Unix
	return homeUnix()
}

func homeUnix() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}
	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}
	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", fmt.Errorf("blank output when reading home directory")
	}
	return result, nil
}

func homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", fmt.Errorf("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}
	return home, nil
}

//CreateMultiDir 调用os.MkdirAll递归创建文件夹
func CreateMultiDir(filePath string) error {
	if pathIsExist(filePath) {
		return nil
	}
	return os.MkdirAll(filePath, os.ModePerm)
}

// 判断所给路径文件/文件夹是否存在(返回true是存在)
func pathIsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
