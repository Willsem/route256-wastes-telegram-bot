package enums

type UserContext int

const (
	NoContext UserContext = iota
	AddWaste
	ChangeCurrency
	SetLimit
)
