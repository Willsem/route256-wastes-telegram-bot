package enums

import "fmt"

type CommandType string

const (
	CommandTypeAdd         CommandType = "/add"
	CommandTypeSetLimit    CommandType = "/setLimit"
	CommandTypeGetLimit    CommandType = "/getLimit"
	CommandTypeWeekReport  CommandType = "/week"
	CommandTypeMonthReport CommandType = "/month"
	CommandTypeYearReport  CommandType = "/year"
	CommandTypeCurrency    CommandType = "/currency"

	CommandTypeUnknown CommandType = ""
)

func ParseCommandType(command string) (CommandType, error) {
	switch command {
	case string(CommandTypeAdd):
		return CommandTypeAdd, nil
	case string(CommandTypeSetLimit):
		return CommandTypeSetLimit, nil
	case string(CommandTypeGetLimit):
		return CommandTypeGetLimit, nil
	case string(CommandTypeWeekReport):
		return CommandTypeWeekReport, nil
	case string(CommandTypeMonthReport):
		return CommandTypeMonthReport, nil
	case string(CommandTypeYearReport):
		return CommandTypeYearReport, nil
	case string(CommandTypeCurrency):
		return CommandTypeCurrency, nil
	default:
		return CommandTypeUnknown, fmt.Errorf("Unknown command type")
	}
}
