package converter

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_workflowIdentifierToFileName(t *testing.T) {
	tests := []struct {
		input string
		want  string
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
			"hi😁",
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("converting `%s'", tt.input), func(t *testing.T) {
			assert.Equal(t, tt.want, workflowIdentifierToFileName(tt.input))
		})
	}
}

func Test_toID(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			"some spaces  here",
			"someSpacesHere",
		},
		{
			"@thing/other workflow",
			"build",
		},
		{
			"weirdly\x00broken\x00 \a\a\a",
			"build",
		},
		{
			"------",
			"build",
		},
		{
			"hi😁",
			"build",
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("converting `%s'", tt.input), func(t *testing.T) {
			assert.Equal(t, tt.want, toID(tt.input))
		})
	}
}
