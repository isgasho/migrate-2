package converter

import (
	"fmt"
	"io"
	"strings"

	"github.com/actions/workflow-parser/model"
	"github.com/actions/workflow-parser/parser"
	"github.com/github/mu/errs"
	"gopkg.in/yaml.v2"
)

const workflowDirectory = ".github/workflows"

// Parse accepts V1 ACL and converts it into a data-structure representing the V2 format
func Parse(v1Workflow io.Reader) (*parsed, error) {
	actions, err := parser.Parse(v1Workflow)
	if err != nil {
		return nil, errs.Wrap(err, "Invalid workflow file")
	}

	converted, err := fromConfiguration(actions)
	if err != nil {
		return nil, err
	}

	return converted, nil
}

// fromConfiguration takes the v1 AST and converts into our v2 data-structure
func fromConfiguration(configuration *model.Configuration) (*parsed, error) {
	converted := parsed{}

	actByID := make(map[string]*model.Action, len(configuration.Actions))
	for _, act := range configuration.Actions {
		actByID[act.Identifier] = act
	}

	fn := newFilenames(configuration.Workflows)

	for i, wf := range configuration.Workflows {
		w := workflow{
			Name: wf.Identifier,
			Jobs: make(map[string]job, 0),
		}
		writeOn(&w, wf.On)
		// Make a job per resolve target
		acts, err := serializeWorkflow(wf, actByID)
		if err != nil {
			return nil, err
		}

		j := job{
			RunsOn: "ubuntu-latest",
		}
		resolved := acts[0]
		id := toID(resolved.Identifier)
		if id != resolved.Identifier {
			j.Name = resolved.Identifier
		}

		j.Actions = make([]*action, 0, len(acts))
		for _, a := range acts {
			ca := &action{
				Uses: a.Uses.String(),
				Name: a.Identifier,
				Env:  a.Env,
			}
			if a.Runs != nil {
				ca.Entrypoint = convertCommandExpressions(a.Runs.Split())
			}
			if a.Args != nil {
				ca.Args = convertCommandExpressions(a.Args.Split())
			}
			j.Actions = append(j.Actions, ca)
		}
		w.Jobs[id] = j

		// if we have a single workflow for an event, name the file after that event
		w.fileName = fn.create(wf, i)

		converted.workflows = append(converted.workflows, &w)
	}

	return &converted, nil
}

func convertCommandExpressions(ss []string) string {
	out := make([]string, 0, len(ss))
	for _, s := range ss {
		out = append(out, convertGithubEnvironmentReferences(s))
	}
	return strings.Join(out, " ")
}

var replacements = []string{
	"$GITHUB_WORKFLOW", "${{ github.workflow }}",
	"$GITHUB_ACTION", "${{ github.action }}",
	"$GITHUB_ACTOR", "${{ github.actor }}",
	"$GITHUB_REPOSITORY", "${{ github.repository }}",
	"$GITHUB_EVENT_NAME", "${{ github.event_name }}",
	"$GITHUB_EVENT_PATH", "${{ github.event_path }}",
	"$GITHUB_WORKSPACE", "${{ github.workspace }}",
	"$GITHUB_SHA", "${{ github.sha }}",
	"$GITHUB_REF", "${{ github.ref }}",
	"$GITHUB_TOKEN", "${{ github.token }}",
}
var envVarReplacements = strings.NewReplacer(replacements...)

func convertGithubEnvironmentReferences(s string) string {
	return envVarReplacements.Replace(s)
}

func writeOn(w *workflow, on model.On) {
	if o, ok := on.(*model.OnSchedule); ok {
		w.OnSchedule = map[string][]map[string]string{
			"schedules": {
				{
					"cron": o.Expression,
				},
			},
		}
	} else {
		w.On = on.String()
	}
}

func onToEvent(on model.On) string {
	if _, ok := on.(*model.OnSchedule); ok {
		return "schedule"
	}
	return on.String()
}

// serializeWorkflow takes the graph from a v1 workflow, and serializes it via a breadth-first-search (which
// topographically sorts any valid workflow).
func serializeWorkflow(workflow *model.Workflow, actByID map[string]*model.Action) ([]*model.Action, error) {
	reverseRoute := make([]*model.Action, 0)
	seen := make(map[string]struct{})

	queue := make([]*model.Action, 0)
	for _, resolveID := range workflow.Resolves {
		act, ok := actByID[resolveID]
		if !ok {
			return nil, errs.Errorf("Resolves to invalid action `%s'", resolveID)
		}
		queue = append(queue, act)
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if _, ok := seen[current.Identifier]; !ok {
			reverseRoute = append(reverseRoute, current)
			seen[current.Identifier] = struct{}{}
		}

		for _, needed := range current.Needs {
			act, ok := actByID[needed]
			if !ok {
				return nil, errs.Errorf("Resolves to invalid action `%s'", needed)
			}
			queue = append(queue, act)
		}
	}

	l := len(reverseRoute)
	plan := make([]*model.Action, l)
	for i, n := range reverseRoute {
		plan[l-i-1] = n
	}

	return plan, nil
}

// Files returns the set of v2 workflow files required to perform the work specified in the V1 ACL file
func (converted *parsed) Files() ([]OutputFile, error) {
	of := make([]OutputFile, 0, len(converted.workflows))
	for _, wf := range converted.workflows {
		j, err := yaml.Marshal(wf)
		if err != nil {
			return nil, err
		}

		y := string(j)
		y = strings.Replace(y, "__workflowKeyOn__", "on", 1)
		y = strings.Replace(y, "__workflowKeyOnSchedule__", "on", 1)

		of = append(of, OutputFile{
			Path:    fmt.Sprintf("%s/%s", workflowDirectory, wf.fileName),
			Content: string(y),
		})
	}

	return of, nil
}
