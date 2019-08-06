package converter

import (
	"fmt"
	"regexp"
	"strings"
)

var spaces = regexp.MustCompile(`\s+`)

var safeID = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)

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
	// be pretty conservative - we only have one job so there's no clashng, and we set name if the build isn't equal to
	// identifer anyway, so users' names will be used.
	c := op.String()
	if safeID.MatchString(c) {
		return c
	}
	return "build"
}

var notUnicodeLetterRE = regexp.MustCompile(`[^\pL\pS\pN]+`)

func workflowIdentifierToFileName(s string) string {
	s = strings.TrimSpace(s)
	s = notUnicodeLetterRE.ReplaceAllString(s, "-")
	substrs := spaces.Split(strings.ToLower(s), -1)
	s = strings.Join(substrs, "-")
	return strings.Trim(s, "-")
}
