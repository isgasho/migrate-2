package migrate

type action struct {
	Name string `yaml:",omitempty"`
	Uses string
	Args string `yaml:",omitempty"`
}

// converted over from a workflow
type job struct {
	Name    string    `yaml:",omitempty"`
	Actions []*action `yaml:"steps"`
}

type workflow struct {
	On       string
	Name     string
	fileName string
	Jobs     map[string]job
}

type parsed struct {
	workflows []*workflow
}

type OutputFile struct {
	Path    string
	Content string
}
