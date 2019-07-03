package parser

func isValidOn(on string) bool {
	return IsSchedule(on) || isAllowedEventType(on)
}
