package bitopro

import (
	"encoding/json"
	"time"

	"github.com/bitoex/bitopro-api-go/pkg/ws"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	_restfulHost = "https://api.bitopro.com/v3"
	_wsHost      = "wss://stream.bitopro.com:443/ws"
)

func ConnectPublicWs() (*ws.Ws, error) {
	wss := ws.NewPublicWs()
	if wss == nil {
		return nil, errors.New("connect to public ws")
	}
	return wss, nil
}

func ConnectPrivateWs() (*ws.Ws, error) {
	wss := ws.NewPrivateWs(viper.GetString("api.bitopro.email"), viper.GetString("api.bitopro.key"), viper.GetString("api.bitopro.secret"))
	if wss == nil {
		return nil, errors.New("connect to private ws")
	}

	return wss, nil
}

func GeneralPayload() string {
	data := map[string]any{
		"identity": viper.GetString("api.bitopro.email"),
		"nonce":    time.Now().UnixMicro(),
	}
	b, _ := json.Marshal(data)
	return string(b)
}

func CreateOrderPayload() string {
	data := map[string]any{
		"action":    "BUY",
		"type":      "limit",
		"price":     "1.123456789",
		"amount":    "666",
		"timestamp": 0,
	}
	b, _ := json.Marshal(data)
	return string(b)
}
