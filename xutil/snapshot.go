package xutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const snapshotTempExt = ".temp"

func StoreSnapshot(path string, data interface{}) error {
	err := MkdirIfNotExist(filepath.Dir(path))
	if err != nil {
		return fmt.Errorf("store data failed when ensure exist dir, path: %s, error: %v", path, err)
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("store data failed when marshal data, path: %s, error: %v", path, err)
	}

	err = ioutil.WriteFile(path+snapshotTempExt, bytes, 0644)
	if err != nil {
		return fmt.Errorf("store data failed when write file, path: %s, error: %v", path, err)
	}

	err = os.Rename(path+snapshotTempExt, path)
	if err != nil {
		return fmt.Errorf("store data failed when rename file, path: %s, error: %v", path, err)
	}

	return nil
}

func LoadSnapshot(path string, data interface{}) error {
	ok, err := IsFileExist(path)
	if err != nil {
		return fmt.Errorf("load data failed when check file, path: %s, error: %v", path, err)
	}
	if !ok {
		return nil
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("load data failed when read file, path: %s, error: %v", path, err)
	}

	return json.Unmarshal(bytes, data)
}
