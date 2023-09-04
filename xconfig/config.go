package xconfig

import (
	"fmt"

	"github.com/spf13/viper"
)

func Init(path string, c interface{}) error {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("read config error [%v]", err)
	}
	if err := viper.Unmarshal(c); err != nil {
		return fmt.Errorf("unmarshal config json error [%v]", err)
	}

	return nil
}
