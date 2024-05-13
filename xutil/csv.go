package xutil

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/yyliziqiu/xlib/xfile"
	"github.com/yyliziqiu/xlib/xif"
)

func SaveCSV(filename string, rows [][]string) error {
	// 创建存储目录
	err := xfile.MkdirIfNotExist(filepath.Dir(filename))
	if err != nil {
		return fmt.Errorf("mkdir failed [%v]", err)
	}

	// 创建 CSV 文件
	if !strings.HasSuffix(filename, ".csv") {
		filename = filename + ".csv"
	}
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create CSV failed [%v]", err)
	}
	defer file.Close()

	// 写入 CSV 文件
	writer := csv.NewWriter(file)
	defer writer.Flush()
	err = writer.WriteAll(rows)
	if err != nil {
		return fmt.Errorf("write CSV failed [%v]", err)
	}

	return nil
}

func SaveCSV2(filename string, models []any) error {
	if len(models) == 0 {
		return nil
	}
	rows := make([][]string, 0, len(models)+1)
	rows = append(rows, modelFields(models[0]))
	for _, model := range models {
		rows = append(rows, modelValues(model))
	}
	return SaveCSV(filename, rows)
}

func modelFields(model any) []string {
	mt := reflect.TypeOf(model)
	var fields []string
	for i := 0; i < mt.NumField(); i++ {
		f := mt.Field(i)
		header := f.Tag.Get("csv")
		fields = append(fields, xif.If(header != "", header, f.Name))
	}
	return fields
}

func modelValues(model any) []string {
	mv := reflect.ValueOf(model)
	var values []string
	for i := 0; i < mv.NumField(); i++ {
		values = append(values, fmt.Sprintf("%v", mv.Field(i).Interface()))
	}
	return values
}
