package converter

type action struct {
	Name       string `yaml:",omitempty"`
	Uses       string
	Args       string            `yaml:",omitempty"`
	Env        map[string]string `yaml:",omitempty"`
	Entrypoint string            `yaml:",omitempty"`
}

// converted over from a workflow
type job struct {
	Name    string    `yaml:",omitempty"`
	RunsOn     string `yaml:"runs-on"`
	Actions []*action `yaml:"steps"`
}

type workflow struct {
	On         string                         `yaml:"__workflowKeyOn__,omitempty"`
	OnSchedule map[string][]map[string]string `yaml:"__workflowKeyOnSchedule__,omitempty"`
	Name       string
	fileName   string
	Jobs       map[string]job
}

type parsed struct {
	workflows []*workflow
}

type OutputFile struct {
	Path    string
	Content string
}
