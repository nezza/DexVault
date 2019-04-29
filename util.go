package main

import (
	"github.com/binance-chain/go-sdk/common"
)

func (c *CreateOrder) CombinedSymbol() string {
	return common.CombineSymbol(c.BaseAssetSymbol, c.QuoteAssetSymbol)
}

func (c *CancelOrder) CombinedSymbol() string {
	return common.CombineSymbol(c.BaseAssetSymbol, c.QuoteAssetSymbol)
}
