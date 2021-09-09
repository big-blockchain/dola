// The toolbelt is a set of helper functions that ease the cross usage of strategies.
package dola

import (
	"errors"

	"github.com/thrasher-corp/gocryptotrader/exchanges/account"
	"github.com/thrasher-corp/gocryptotrader/exchanges/ticker"
)

var (
	ErrNeedBalancesStrategy = errors.New("Keep should be configured with balances support")
	ErrCast                 = errors.New("casting failed")
)

func CurrencyBalance(k *Keep, exchangeName, currencyCode string, accountIndex int) (account.Balance, error) {
	st, err := k.Root.Get("balances")
	if errors.Is(err, ErrStrategyNotFound) {
		var empty account.Balance

		return empty, ErrNeedBalancesStrategy
	}

	balances, ok := st.(*BalancesStrategy)
	if !ok {
		var empty account.Balance

		return empty, ErrCast
	}

	return balances.Currency(exchangeName, currencyCode, accountIndex)
}

// Ticker casts a void* to ticker.Price.
func Ticker(p interface{}) ticker.Price {
	x, ok := p.(ticker.Price)
	if !ok {
		panic("")
	}

	return x
}