package strategy

import (
	"errors"

	"main/internal/domain"
)

type StdBotProvider struct {
	Kline domain.KlineProvideServer
	Trade domain.TradeServer
}

func (p *StdBotProvider) Validate() error {
	if p.Kline == nil {
		return errors.New("nil kline provider")
	}

	if p.Trade == nil {
		return errors.New("nil trade server")
	}

	return nil
}
