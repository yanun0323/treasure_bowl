package enum

import "github.com/spf13/viper"

const (
	Secondly        CronSpec = "* * * * * *"
	Minutely        CronSpec = "0 * * * * *"
	FifteenMinutely CronSpec = "0 */15 * * * *"
	ThirtyMinutely  CronSpec = "0 */30 * * * *"
	HalfHourly      CronSpec = "0 */30 * * * *"
	Hourly          CronSpec = "0 0 * * * *"
	Daily           CronSpec = "0 0 0 * * *"
)

type CronSpec string

func (s CronSpec) String() string {
	return string(s)
}

func Config() string {
	c := viper.GetString("kline.cron.frequency")
	if len(c) != 0 {
		return c
	}
	return "* * * * * *"
}
