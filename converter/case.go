package converter

import (
	"fmt"
	"regexp"
	"strings"
)

var spaces = regexp.MustCompile(`\s+`)

func toID(s string) string {
	s = strings.TrimSpace(s)
	substrs := spaces.Split(s, -1)
	op := strings.Builder{}
	for i, sub := range substrs {
		if len(sub) <= 1 {
			op.WriteString(sub)
			continue
		}

		if i == 0 {
			sub = fmt.Sprintf("%s%s",
				strings.ToLower(string(sub[0])),
				sub[1:],
			)
		} else {
			sub = fmt.Sprintf("%s%s",
				strings.ToUpper(string(sub[0])),
				sub[1:],
			)
		}
		op.WriteString(sub)
	}
	return op.String()
}

var confusingFileNameRE = regexp.MustCompile(`[@\!{}$/\s]+`)

func workflowIdentifierToFileName(s string) string {
	s = strings.TrimSpace(s)
	s = confusingFileNameRE.ReplaceAllString(s, "-")
	substrs := spaces.Split(strings.ToLower(s), -1)
	return strings.Join(substrs, "-")
}
