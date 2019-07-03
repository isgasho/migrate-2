package parser

import (
	"regexp"
)

var scheduleRegex = regexp.MustCompile(`\Aschedule\(([\*\s\/\d,]+)\)\z`)

func IsSchedule(onString string) bool {
	return scheduleRegex.MatchString(onString)
}

func isScheduleExpression(expression string) bool {
	return extractScheduleExpression(expression) == ""
}

func extractScheduleExpression(onString string) string {
	matches := scheduleRegex.FindAllStringSubmatch(onString, 1)
	if len(matches) == 0 {
		return ""
	}
	return matches[0][1]
}
