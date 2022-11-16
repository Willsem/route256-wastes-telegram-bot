package requests

type Period int

const (
	PeriodWeek Period = iota
	PeriodMonth
	PeriodYear
)

//easyjson:json
type GetReport struct {
	UserID              int64   `json:"user_id"`
	Period              Period  `json:"period"`
	CurrencyExchange    float64 `json:"currency_exchange"`
	CurrencyDesignation string  `json:"currency_designation"`
}
