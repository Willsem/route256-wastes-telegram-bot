package models

type ExchangeData struct {
	Base  string
	Rates map[string]float64
}

func NewExchangeData(base string, rates map[string]float64) *ExchangeData {
	return &ExchangeData{
		Base:  base,
		Rates: rates,
	}
}
