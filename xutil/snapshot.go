package xutil

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const snapshotTempExt = ".temp"

func NewSnapshot(path string) *Snapshot {
	return &Snapshot{Path: path}
}

type Snapshot struct {
	Path string
}

func (s *Snapshot) Store(data interface{}) error {
	err := MkdirIfNotExist(filepath.Dir(s.Path))
	if err != nil {
		return fmt.Errorf("store data failed when ensure exist dir, path: %s, error: %v", s.Path, err)
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("store data failed when marshal data, path: %s, error: %v", s.Path, err)
	}

	err = os.WriteFile(s.Path+snapshotTempExt, bytes, 0644)
	if err != nil {
		return fmt.Errorf("store data failed when write file, path: %s, error: %v", s.Path, err)
	}

	err = os.Rename(s.Path+snapshotTempExt, s.Path)
	if err != nil {
		return fmt.Errorf("store data failed when rename file, path: %s, error: %v", s.Path, err)
	}

	return nil
}

func (s *Snapshot) Load(data interface{}) error {
	ok, err := IsFileExist(s.Path)
	if err != nil {
		return fmt.Errorf("load data failed when check file, path: %s, error: %v", s.Path, err)
	}
	if !ok {
		return nil
	}

	bytes, err := os.ReadFile(s.Path)
	if err != nil {
		return fmt.Errorf("load data failed when read file, path: %s, error: %v", s.Path, err)
	}

	return json.Unmarshal(bytes, data)
}
