package xfile

import (
	"os"
)

func Exist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func MkdirIfNotExist(path string) error {
	exist, err := Exist(path)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	return os.MkdirAll(path, 0755)
}
