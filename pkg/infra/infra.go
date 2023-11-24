package infra

import "github.com/yanun0323/pkg/config"

func Init(cfgName string) error {
	return config.Init(cfgName, true, "../config", "../../config", "../../../config", "./config")
}
