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
	Data interface{}
}

func NewSnapshot(path string, data interface{}) *Snapshot {
	return &Snapshot{
		Path: path,
		Data: data,
	}
}

func (s *Snapshot) Save() error {
	return s.SaveData(s.Data)
}

func (s *Snapshot) SaveData(data interface{}) error {
	err := MkdirIfNotExist(filepath.Dir(s.Path))
	if err != nil {
		return fmt.Errorf("mkdir snapshot dir [%s] failed [%v]", filepath.Dir(s.Path), err)
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal snapshot data [%s] failed [%v]", s.Path, err)
	}

	err = os.WriteFile(s.Path+snapshotTempExt, bytes, 0644)
	if err != nil {
		return fmt.Errorf("save snapshot file [%s] failed [%v]", s.Path, err)
	}

	err = os.Rename(s.Path+snapshotTempExt, s.Path)
	if err != nil {
		return fmt.Errorf("rename snapshot file [%s] failed [%v]", s.Path, err)
	}

	return nil
}

func (s *Snapshot) Load() error {
	return s.LoadData(s.Data)
}

func (s *Snapshot) LoadData(data interface{}) error {
	ok, err := IsFileExist(s.Path)
	if err != nil {
		return fmt.Errorf("check snapshot file [%s] failed [%v]", s.Path, err)
	}
	if !ok {
		return nil
	}

	bytes, err := os.ReadFile(s.Path)
	if err != nil {
		return fmt.Errorf("load snapshot file [%s] failed [%v]", s.Path, err)
	}

	return json.Unmarshal(bytes, data)
}
