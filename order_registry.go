package dola

import (
	"sync"

	"github.com/thrasher-corp/gocryptotrader/exchanges/order"
)

type OrderKey struct {
	ExchangeName string
	OrderID      string
}

type OrderValue struct {
	SubmitResponse order.SubmitResponse
	UserData       interface{}
}

type OrderRegistry struct {
	m sync.Map
}

func NewOrderRegistry() OrderRegistry {
	return OrderRegistry{sync.Map{}}
}

func (r *OrderRegistry) OnSubmit(exchangeName string, response order.SubmitResponse, userData interface{}) {
	key := OrderKey{
		ExchangeName: exchangeName,
		OrderID:      response.OrderID,
	}
	value := OrderValue{
		SubmitResponse: response,
		UserData:       userData,
	}
	r.m.Store(key, value)
}
