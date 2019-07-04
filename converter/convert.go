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

	countsByEvent := make(map[string]int)
	for _, wf := range configuration.Workflows {
		event := onToEvent(wf.On)
		countsByEvent[event] = countsByEvent[event] + 1
	}

	fileNames := make(map[string]struct{})

	for i, wf := range configuration.Workflows {
		// TODO schedules
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

		j := job{}
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
				ca.Entrypoint = strings.Join(a.Runs.Split(), " ")
			}
			if a.Args != nil {
				ca.Args = strings.Join(a.Args.Split(), " ")
			}
			j.Actions = append(j.Actions, ca)
		}
		w.Jobs[id] = j

		// if we have a single workflow for an event, name the file after that event
		ev := onToEvent(wf.On)
		if countsByEvent[ev] == 1 {
			w.fileName = fmt.Sprintf("%s.yml", ev)
		} else {
			converted := workflowIdentifierToFileName(wf.Identifier)
			// if identifier can't be converted to something meaningful, just use a number
			if converted == "" {
				converted = fmt.Sprintf("%v", i + 1)
			}
			// or if due to conversion we end up with a duplicate, use suffix to make unique
			if _, ok := fileNames[converted]; ok {
				converted = fmt.Sprintf("%s-%v", converted, i + 1)
			}
			fileNames[converted] = struct{}{}
			w.fileName = fmt.Sprintf("%s-%s.yml", ev, converted)
		}

		converted.workflows = append(converted.workflows, &w)
	}

	return &converted, nil
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
