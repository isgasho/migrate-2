package converter

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// workflowToYAML outputs a YAML string from a workflow, taking advantage of the
// specific and limited input of our workflows to keep the implementation simple.
func workflowToYAML(wf *workflow) string {
	if len(wf.Jobs) > 1 {
		panic("Invariant: assuming conversion outputs max one job")
	}

	buf := &strings.Builder{}
	w := &stackedWriter{context: &rootWriter{buf}}

	outputComments(w, wf.comments)
	writeLineF(w, "on: %s", yamlQuote(wf.On))
	writeLineF(w, "name: %s", yamlQuote(wf.Name))
	writeLineF(w, "\njobs:")
	w = w.push(&mapContext{})

	for _, job := range wf.Jobs {
		writeLineF(w, "%s:", toID(job.Name))
		w = w.push(&mapContext{})
		if job.Name != "" {
			writeLineF(w, "name: %s", yamlQuote(job.Name))
		}
		w.WriteLine("steps:")

		list := &listContext{}
		w = w.push(list)
		for _, act := range job.Actions {
			if act.Name == "" {
				writeLineF(w, "uses: %s", yamlQuote(act.Uses))
			} else {
				writeLineF(w, "name: %s", yamlQuote(act.Name))
				writeLineF(w, "uses: %s", yamlQuote(act.Uses))
			}
			if act.Args != "" {
				writeLineF(w, "args: %s", yamlQuote(act.Args))
			}
			list.itemEnd()
		}
		w = w.Pop() // list
		w = w.Pop() // map
	}

	return buf.String()
}

func writeLineF(w *stackedWriter, f string, vars ...interface{}) {
	w.WriteLine(fmt.Sprintf(f, vars...))
}

var safeStartRE = regexp.MustCompile(`^[a-zA-Z]`)

func yamlQuote(s string) string {
	if safeStartRE.MatchString(s) && !strings.Contains(s, `"`) {
		return s
	}
	// TODO
	return fmt.Sprintf("%q", strconv.Quote(s))
}

func outputComments(w lineWriter, cs comments) {
	if len(cs) == 0 {
		return
	}
	// TODO proper comment handling
	for _, cg := range cs {
		txt := make([]string, 0, len(cs))
		for _, n := range cg.List {
			txt = append(txt, n.Text)
		}
		w.WriteLine(fmt.Sprintf("# %s", strings.Join(txt, " ")))
	}
}

type listContext struct {
	afterFirst bool
}

func (lc *listContext) WriteLine(line string) string {
	if lc.afterFirst {
		return fmt.Sprintf("  %s", line)
	} else {
		lc.afterFirst = true
		return fmt.Sprintf("- %s", line)
	}
}

func (lc *listContext) itemEnd() {
	lc.afterFirst = false
}

type lineWriter interface {
	WriteLine(s string) string
}

type mapContext struct{}

func (mc *mapContext) WriteLine(s string) string {
	return fmt.Sprintf("  %s", s)
}

// implements a stack of writers to handle printing of hierarchical structures
type stackedWriter struct {
	parent  *stackedWriter
	context lineWriter
}

func (sw *stackedWriter) WriteLine(op string) string {
	n := sw
	for n != nil {
		op = n.context.WriteLine(op)
		n = n.parent
	}
	return op
}

func (sw *stackedWriter) Pop() *stackedWriter {
	if sw.parent == nil {
		panic("popped too much, stop popping!")
	}
	return sw.parent
}
func (sw *stackedWriter) push(context lineWriter) *stackedWriter {
	return &stackedWriter{
		context: context,
		parent:  sw,
	}
}

type rootWriter struct {
	w io.Writer
}

func (rw *rootWriter) WriteLine(s string) string {
	rw.w.Write([]byte(s))
	rw.w.Write([]byte("\n"))
	return ""
}
