package util

import (
	"github.com/spf13/viper"
	"github.com/yanun0323/pkg/logs"
)

func LogLevel() logs.Level {
	return logs.NewLevel(viper.GetString("log.level"))
}
