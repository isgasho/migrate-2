package converter

import (
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
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@master
    - name: action one
      uses: docker://alpine
      with:
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
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@master
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

func TestConvertOneSchedules(t *testing.T) {
	assertCorrect(t, eg{
		input: `workflow "scheduled" {
  on = "schedule(* * * * *)"
  resolves = [
    "A",
  ]
}

workflow "scheduled two" {
  on = "schedule(* * * * *)"
  resolves = [
    "A",
  ]
}

action "A" { 
  uses = "./A" 
}
`,
		output: map[string]string{
			".github/workflows/schedule-scheduled.yml": `on:
  schedules:
  - cron: '* * * * *'
name: scheduled
jobs:
  A:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@master
    - name: A
      uses: ./A
`,
			".github/workflows/schedule-scheduled-two.yml": `on:
  schedules:
  - cron: '* * * * *'
name: scheduled two
jobs:
  A:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@master
    - name: A
      uses: ./A
`,
		},
	})
}

func TestMultipleWorkflows(t *testing.T) {
	assertCorrect(t, eg{
		input: `workflow "one" {
  on = "push"
  resolves = [
    "A",
  ]
}

workflow "B" {
  on = "push"
  resolves = [
    "A",
  ]
}

workflow "C" {
  on = "repository_dispatch"
  resolves = [
    "A",
  ]
}

action "A" { 
  uses = "./A" 
}

`,
		output: map[string]string{
			".github/workflows/push-one.yml": `on: push
name: one
jobs:
  A:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@master
    - name: A
      uses: ./A
`,
			".github/workflows/push-b.yml": `on: push
name: B
jobs:
  A:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@master
    - name: A
      uses: ./A
`,
			".github/workflows/repository_dispatch.yml": `on: repository_dispatch
name: C
jobs:
  A:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@master
    - name: A
      uses: ./A
`,
		},
	})
}

func TestInvalidResolves(t *testing.T) {
	_, err := Parse(strings.NewReader(`
workflow "a" {
    on = "push"
	resolves = ["not there"]
}
`))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "resolves unknown action")
}

func TestEnv(t *testing.T) {
	assertCorrect(t, eg{
		input: `workflow "one" {
  on = "push"
  resolves = [
    "A",
  ]
}

action "A" { 
  uses = "./A" 
  env = {
    ONE = "one v"
    "TWO" = "two v"
  }
}
`,
		output: map[string]string{
			".github/workflows/push.yml": `on: push
name: one
jobs:
  A:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@master
    - name: A
      uses: ./A
      env:
        ONE: one v
        TWO: two v
`,
		},
	})
}

func TestRuns(t *testing.T) {
	assertCorrect(t, eg{
		input: `workflow "workflow one" {
  on = "push"
  resolves = [
    "action one",
  ]
}

action "action one" {
  uses = "docker://alpine"
  runs = ["sh", "-c", "echo hi there"] 
}`,
		output: map[string]string{
			".github/workflows/push.yml": `on: push
name: workflow one
jobs:
  actionOne:
    name: action one
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@master
    - name: action one
      uses: docker://alpine
      with:
        entrypoint: sh
        args: -c "echo hi there"
`,
		},
	})
}

func TestGithubEnvironmentVariableRewriting(t *testing.T) {
	assertCorrect(t, eg{
		input: `workflow "workflow one" {
  on = "push"
  resolves = [
    "action one",
  ]
}

action "action one" {
  uses = "docker://alpine"
  runs = ["sh", "-c", "echo $GITHUB_SHA"] 
}`,
		output: map[string]string{
			".github/workflows/push.yml": `on: push
name: workflow one
jobs:
  actionOne:
    name: action one
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@master
    - name: action one
      uses: docker://alpine
      with:
        entrypoint: sh
        args: -c "echo ${{ github.sha }}"
`,
		},
	})
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
			assert.Equal(t, op, f.Content)
		} else {
			assert.Failf(t, "unexpected path %s", f.Path)
		}
	}
	assert.ElementsMatch(t, expectedFiles, allFiles)
}
