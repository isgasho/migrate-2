package converter

import (
	"fmt"
	"github.com/actions/workflow-parser/model"
)

type fileNameSet struct {
	countsByEvent map[string]int
	fileNames     map[string]struct{}
}

func newFilenames(wfs []*model.Workflow) *fileNameSet {
	fs := fileNameSet{
		fileNames: make(map[string]struct{}),
	}
	countsByEvent := make(map[string]int)
	for _, wf := range wfs {
		event := onToEvent(wf.On)
		countsByEvent[event] = countsByEvent[event] + 1
	}
	fs.countsByEvent = countsByEvent
	return &fs
}

func (fs *fileNameSet ) create(wf *model.Workflow, i int) string {
	ev := onToEvent(wf.On)
	fn := ""
	if fs.countsByEvent[ev] == 1 {
		fn = fmt.Sprintf("%s.yml", ev)
	} else {
		converted := workflowIdentifierToFileName(wf.Identifier)
		// if identifier can't be converted to something meaningful, just use a number
		if converted == "" {
			converted = fmt.Sprintf("%v", i + 1)
		}
		// if we collapsed to something non-unique, append index
		if _, ok := fs.fileNames[converted]; !ok  {
			converted = fmt.Sprintf("%s-%v", converted, i + 1)
		}
		fn = fmt.Sprintf("%s-%s.yml", ev, converted)
	}
	fs.fileNames[fn] = struct{}{}
	return fn
}
