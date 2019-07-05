package converter

import (
	"fmt"
	"github.com/actions/workflow-parser/model"
)

const itAintWorthHangingOverThis = 10000

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

func (fs *fileNameSet) create(wf *model.Workflow, i int) string {
	ev := onToEvent(wf.On)
	fn := ""
	if fs.countsByEvent[ev] == 1 {
		fn = fmt.Sprintf("%s.yml", ev)
	} else {
		id := workflowIdentifierToFileName(wf.Identifier)
		// if identifier can't be converted to something meaningful, just use a number
		// if we collapsed to something non-unique, append index
		for i = 0; i < itAintWorthHangingOverThis; i++ {
			c := id
			if i > 0 {
				c = fmt.Sprintf("%s-%v", id, i+1)
			}
			if _, ok := fs.fileNames[forEvent(ev, c)]; !ok {
				id = c
				break
			}
		}
		fn = forEvent(ev, id)
	}
	fs.fileNames[fn] = struct{}{}
	return fn
}

func forEvent(event, name string) string {
	return fmt.Sprintf("%s-%s.yml", event, name)
}
