package util

import (
	"github.com/spf13/viper"
	"github.com/yanun0323/pkg/logs"
)

func LogLevel() uint16 {
	level, err := logs.NewLevel(viper.GetString("log.level"))
	if err != nil {
		return 0
	}
	return uint16(level)
}
