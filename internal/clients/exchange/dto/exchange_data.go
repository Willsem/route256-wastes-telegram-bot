package dto

type ExchangeData struct {
	Base    string             `json:"base"`
	Rates   map[string]float64 `json:"rates"`
	Success bool               `json:"success"`
}
