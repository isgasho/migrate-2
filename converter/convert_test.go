package converter

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type eg struct {
	input  string
	output map[string]string
}

func TestSimpleConversion(t *testing.T) {

	assertCorrect(t, eg{
		input: `workflow "workflow one" {
  on = "push"
  resolves = [
    "action one",
  ]
}

action "action one" {
  uses = "docker://alpine"
  args = "echo hi"
}`,
		output: map[string]string{
			".github/workflows/push.yml": `on: push
name: workflow one

jobs:
  actionOne:
    name: action one
    steps:
    - name: action one
      uses: docker://alpine
      args: echo hi
`,
		},
	})
}

func TestConvertParallelWorkflows(t *testing.T) {
	assertCorrect(t, eg{
		input: `workflow "fan in out" {
  on = "push"
  resolves = [
    "E",
    "F",
  ]
}

action "A" { 
  uses = "./A" 
}
action "B" {
  uses = "./A"
  needs = ["A"]
}
action "C" {
  uses = "./A"
  needs = ["A"]
}
action "D" {
  uses = "./A"
  needs = ["B", "C"]
}
action "E" {
  uses = "./A"
  needs = ["D"]
}
action "F" {
  uses = "./A"
  needs = ["D"]
}
`,
		output: map[string]string{
			".github/workflows/push.yml": `on: push
name: fan in out
jobs:
  A:
    steps:
    - name: A
      uses: ./A
    - name: C
      uses: ./A
    - name: B
      uses: ./A
    - name: D
      uses: ./A
    - name: F
      uses: ./A
    - name: E
      uses: ./A
`,
		},
	})
}

func TestComments(t *testing.T) {
	assertCorrect(t, eg{
		input: `
// Some header
/* different comment */

// Workflow comment
workflow "workflow one" {
  // on comment
  on = "push"
  resolves = [
    "action one", // line comment
  ]
}

action "action one" {
 // pre comment
  uses = "docker://alpine" // attribute line comment
}`,
		output: map[string]string{
			".github/workflows/push.yml": `
# Some header
# different comment
# Workflow comment
# on comment
"on": push
name: workflow one
jobs:
  actionOne:
    name: action one
    steps:
    - name: action one
      uses: docker://alpine # attribute line comment
      args: echo hi
`,
		},
	})
}

func TestSchedules(t *testing.T) {
	t.Skip("TODO")
}

func TestMultipleWorkflows(t *testing.T) {
	t.Skip("TODO")
}

func TestInvalidResolves(t *testing.T) {
	t.Skip("TODO")
}

func TestInvalidNeeds(t *testing.T) {
	t.Skip("TODO")
}

func TestEnv(t *testing.T) {
	t.Skip("TODO")
}

func assertCorrect(t *testing.T, eg eg) {
	op, err := Parse(strings.NewReader(eg.input))
	require.NoError(t, err)

	files, err := op.Files()
	require.NoError(t, err)

	expectedFiles := make([]string, 0)
	for p := range eg.output {
		expectedFiles = append(expectedFiles, p)
	}

	allFiles := make([]string, 0)
	for _, f := range files {
		allFiles = append(allFiles, f.Path)
		if op, ok := eg.output[f.Path]; ok {
			matched := assert.Equal(t, op, f.Content)
			if !matched {
				fmt.Println(f.Content)
			}
		} else {
			assert.Failf(t, "unexpected path %s", f.Path)
		}
	}
	assert.ElementsMatch(t, expectedFiles, allFiles)
}
