package io

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

func GetDefaultConfigFile(enc Encoding) (string, error) {
	var err error
	var dir = fmt.Sprintf("%s%c%s", HomeFolder(), os.PathSeparator, ".go-cron")
	if ! FileExists(dir) {
		err = CreateFolder(dir, 0777)
		if err != nil {
			return dir, err
		}
	}
	var file = fmt.Sprintf("%s%c%s.%s", dir, os.PathSeparator, "config", enc.String())
	return file, err
}

func GetCurrentPath() string {
	wd, err := os.Getwd()
	if err != nil {
		exec, err := os.Executable()
		if err != nil {
			return HomeFolder()
		}
		return filepath.Dir(exec)
	}
	return wd
}

func HomeFolder() string {
	usr, err := user.Current()
	if err != nil {
		return os.TempDir()
	}
	return usr.HomeDir
}

func FileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func CreateFilePath(path string, perm os.FileMode) (string, error) {
	var err error
	dir, _ := filepath.Split(path)
	if _, err = os.Stat(dir); err != nil {
		err = os.MkdirAll(dir, perm)
	}
	return dir, err
}


func CreateFolder(dir string, perm os.FileMode) error {
	var err error
	if _, err = os.Stat(dir); err != nil {
		err = os.MkdirAll(dir, perm)
	}
	return err
}