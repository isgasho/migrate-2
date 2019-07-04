package converter

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_workflowIdentifierToFileName(t *testing.T) {
	tests := []struct {
		input string
		want string
	}{
		{
			"some spaces  here",
			"some-spaces-here",
		},
		{
			"@thing/other workflow",
			"thing-other-workflow",
		},
		{
			"weirdly\x00broken\x00 \a\a\a",
			"weirdly-broken",
		},
		{
			"------",
			"",
		},
		{
			"hi😁",
			"hi",
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("converting `%s'", tt.input), func(t *testing.T) {
			assert.Equal(t, tt.want, workflowIdentifierToFileName(tt.input))
		})
	}
}
