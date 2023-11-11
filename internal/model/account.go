package model

import "main/internal/util"

type Account struct {
	Balances util.SyncMap[string, Balance]
}

type Balance struct {
	Available string
	InTrade   string
	Locked    string
}
