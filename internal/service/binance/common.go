package binance

import (
	binance_connector "github.com/binance/binance-connector-go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func NewClient(url string) (*binance_connector.Client, error) {
	key := viper.GetString("api.binance.key")
	secret := viper.GetString("api.binance.secret")
	client := binance_connector.NewClient(key, secret, _klineUrl)
	if client == nil {
		return nil, errors.New("unknown connection error")
	}
	return client, nil
}
