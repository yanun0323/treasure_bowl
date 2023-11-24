package util

import (
	"github.com/spf13/viper"
	"github.com/yanun0323/pkg/logs"
)

func LogLevel() logs.Level {
	return logs.NewLevel(viper.GetString("log.level"))
}

func CronSpec() string {
	c := viper.GetString("kline.cron.frequency")
	if len(c) != 0 {
		return c
	}
	return "* * * * * *"
}
