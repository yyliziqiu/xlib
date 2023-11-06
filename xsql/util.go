package xsql

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
)

func IsNoRowsError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func IsMySQLDuplicateKeyError(err error) bool {
	err2, ok := err.(*mysql.MySQLError)
	if !ok {
		return false
	}
	return err2.Number == 1062
}

func IsPostgresDuplicateKeyError(err error) bool {
	err2, ok := err.(*pq.Error)
	if !ok {
		return false
	}
	return err2.Code == "23505"
}

func JoinIntValue(items []int) string {
	length := len(items)
	if length == 0 {
		return ""
	}
	sb := strings.Builder{}
	for i := 0; i < length-1; i++ {
		sb.WriteString(strconv.Itoa(items[i]))
		sb.WriteString(", ")
	}
	sb.WriteString(strconv.Itoa(items[length-1]))
	return sb.String()
}

func JoinStringValue(items []string) string {
	length := len(items)
	if length == 0 {
		return ""
	}
	return "'" + strings.Join(items, "', '") + "'"
}

func JoinStingValueWithSlashes(items []string) string {
	length := len(items)
	for i := 0; i < length; i++ {
		items[i] = AddSlashes(items[i])
	}
	return JoinStringValue(items)
}

func AddSlashes(str string) string {
	chars := []rune(str)
	temp := make([]rune, 0, len(chars))
	for _, c := range chars {
		if c == '\\' || c == '"' || c == '\'' {
			temp = append(temp, '\\')
			temp = append(temp, c)
		} else {
			temp = append(temp, c)
		}
	}
	return string(temp)
}
