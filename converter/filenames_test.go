package converter

import (
	"fmt"
	"github.com/actions/workflow-parser/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_newFilenames(t *testing.T) {
	type eg struct {
		wf     model.Workflow
		output string
	}

	egs := []eg{
		{
			wf: model.Workflow{
				Identifier: "dupe",
				On:         &model.OnEvent{Event: "push"},
			},
			output: "push-dupe.yml",
		},
		{
			wf: model.Workflow{
				Identifier: "dupe",
				On:         &model.OnEvent{Event: "push"},
			},
			output: "push-dupe-2.yml",
		},
		{
			wf: model.Workflow{
				Identifier: "dupe post collapse",
				On:         &model.OnEvent{Event: "push"},
			},
			output: "push-dupe-post-collapse.yml",
		},
		{
			wf: model.Workflow{
				Identifier: "dupe    post collapse",
				On:         &model.OnEvent{Event: "push"},
			},
			output: "push-dupe-post-collapse-2.yml",
		},
		{
			wf: model.Workflow{
				Identifier: "unique event",
				On:         &model.OnEvent{Event: "issues"},
			},
			output: "issues.yml",
		},
		{
			wf: model.Workflow{
				Identifier: "weird\x00name",
				On:         &model.OnEvent{Event: "push"},
			},
			output: "push-weird-name.yml",
		},
	}

	wfs := make([]*model.Workflow, 0, len(egs))
	for _, eg := range egs {
		wf := eg.wf
		wfs = append(wfs, &wf)
	}

	fn := newFilenames(wfs)

	for i, eg := range egs {
		t.Run(fmt.Sprintf("example %v %s", i, eg.wf.Identifier), func(t *testing.T) {
			assert.Equal(t, eg.output, fn.create(&eg.wf, i))
		})
	}

}
