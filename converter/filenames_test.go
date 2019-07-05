package converter

import (
	"github.com/github/workflow-parser/model"
	"testing"
)

func Test_newFilenames(t *testing.T) {
	wfs := []*model.Workflow{
		{
			Identifier: "dupe",
			On: &model.OnEvent{Event:"push"},
		},
		{
			Identifier: "dupe",
			On: &model.OnEvent{Event:"push"},
		},
		{
			Identifier: "post collapse dupe",
			On: &model.OnEvent{Event:"push"},
		},
		{
			Identifier: "post  collapse  dupe",
			On: &model.OnEvent{Event:"push"},
		},
		{
			Identifier: "ok",
			On: &model.OnEvent{Event:"push"},
		},
		{
			Identifier: "unique event",
			On: &model.OnEvent{Event:"issues"},
		},
	}
	fn := newFilenames(wfs)
}

