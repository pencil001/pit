package util

import (
	"io/ioutil"
	"os"
)

func IsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, err
}

func IsDir(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return f.IsDir(), nil
}

func IsEmptyDir(path string) (bool, error) {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return false, err
	}
	return len(entries) == 0, nil
}

func CreateDir(path string) error {
	if err := os.MkdirAll(path, 0666); err != nil {
		return err
	}
	return nil
}

func CreateFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
}
