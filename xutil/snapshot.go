package xutil

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const snapshotTempExt = ".temp"

type Snapshot struct {
	Path string
}

func NewSnapshot(path string) *Snapshot {
	return &Snapshot{Path: path}
}

func (s *Snapshot) Store(data interface{}) error {
	err := MkdirIfNotExist(filepath.Dir(s.Path))
	if err != nil {
		return fmt.Errorf("mkdir snapshot dir [%s] error [%v]", filepath.Dir(s.Path), err)
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal snapshot data [%s] error [%v]", s.Path, err)
	}

	err = os.WriteFile(s.Path+snapshotTempExt, bytes, 0644)
	if err != nil {
		return fmt.Errorf("store snapshot file [%s] error [%v]", s.Path, err)
	}

	err = os.Rename(s.Path+snapshotTempExt, s.Path)
	if err != nil {
		return fmt.Errorf("rename snapshot file [%s] error [%v]", s.Path, err)
	}

	return nil
}

func (s *Snapshot) Load(data interface{}) error {
	ok, err := IsFileExist(s.Path)
	if err != nil {
		return fmt.Errorf("check snapshot file [%s] error [%v]", s.Path, err)
	}
	if !ok {
		return nil
	}

	bytes, err := os.ReadFile(s.Path)
	if err != nil {
		return fmt.Errorf("load snapshot file [%s] error [%v]", s.Path, err)
	}

	return json.Unmarshal(bytes, data)
}
