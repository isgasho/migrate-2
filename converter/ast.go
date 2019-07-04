package converter

import "github.com/actions/hcl/hcl/ast"

type action struct {
	Name string `yaml:",omitempty"`
	Uses string
	Args string `yaml:",omitempty"`
}

type comments map[string]*ast.CommentGroup

// converted over from a workflow
type job struct {
	Name     string    `yaml:",omitempty"`
	Actions  []*action `yaml:"steps"`
	comments comments
}

type workflow struct {
	On       string
	Name     string
	fileName string
	Jobs     map[string]job
	comments comments
}

type parsed struct {
	workflows []*workflow
}

type OutputFile struct {
	Path    string
	Content string
}
