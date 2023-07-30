package task

import (
	"bytes"
	stdcontext "context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/elvuel/temprender/context"
	"github.com/elvuel/temprender/transport"

	"github.com/Masterminds/sprig/v3"
	"github.com/bmatcuk/doublestar/v4"
	"go.uber.org/zap"
)

/*
	TODO
		log
		kick in go standard packge `context`
*/

type Task struct {
	Logger *zap.SugaredLogger `json:"-"`

	BasePath string `json:"base_path,omitempty"`
	Glob     string `json:"glob,omitempty"`

	LeftDelim  string                 `json:"left_delim,omitempty"`
	RightDelim string                 `json:"right_delim,omitempty"`
	FuncMap    map[string]interface{} `json:"-,omitempty"`

	RenderedEndTag string          `json:"rendered_end_tag"`
	PerformCtx     context.Context `json:"perform_ctx,omitempty"`

	Transporters GroupedTransporters `json:"transporters,omitempty"`

	tmpl                      *template.Template
	allocatedTemplateNameList []string
}

func NewTask() *Task {
	return &Task{
		Glob:           "**/*",
		Logger:         zap.NewExample().Sugar(),
		FuncMap:        sprig.FuncMap(),
		LeftDelim:      "{{",
		RightDelim:     "}}",
		RenderedEndTag: "rdone",

		Transporters: make(GroupedTransporters),
	}
}

func (t *Task) initTemplate() {
	if t.tmpl == nil {
		t.tmpl = template.New("").Delims(t.LeftDelim, t.RightDelim).Funcs(t.FuncMap)
	}
}

func (t *Task) RegisterTemplateFunc(name string, funk interface{}) {
	t.FuncMap[name] = funk
}

func (t *Task) AllocateNewTemplate(name, data string) error {
	t.initTemplate()
	name = strings.Replace(name, t.BasePath+"/", "", 1)

	_, err := t.tmpl.New(name).Parse(data)
	if err != nil {
		return err
	}

	t.AllocatedTemplateNameList(true)
	return nil
}

func (t *Task) LoadTemplates() error {
	if t.Glob == "" {
		t.Glob = "**/*.*"
	}

	files, err := doublestar.FilepathGlob(filepath.Join(t.BasePath, t.Glob))
	if err != nil {
		return err
	}

	t.initTemplate()

	for _, file := range files {
		b, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		err = t.AllocateNewTemplate(file, string(b))
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Task) LoadTemplatesFromMap(kv map[string]string) error {
	t.initTemplate()

	for file, data := range kv {
		err := t.AllocateNewTemplate(file, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Task) AppendTemplate(name, data string) error {
	return t.AllocateNewTemplate(name, data)
}

func (t *Task) AllocatedTemplateNameList(reload bool) []string {
	if reload || t.allocatedTemplateNameList == nil {
		result := make([]string, 0)
		for _, tmpl := range t.tmpl.Templates() {
			result = append(result, tmpl.Name())
		}
		t.allocatedTemplateNameList = result
	}

	return t.allocatedTemplateNameList
}

func (t *Task) AllocatedTemplateExisted(name string) bool {
	t.AllocatedTemplateNameList(false)

	for _, val := range t.allocatedTemplateNameList {
		if val == name {
			return true
		}
	}

	return false
}

func (t *Task) RegisterTransporters(name string, transporters ...transport.Transporter) error {
	if t.Transporters == nil {
		t.Transporters = make(GroupedTransporters)
	}

	if t.Transporters[name] == nil {
		t.Transporters[name] = make(transport.Transporters, 0)
	}

	t.Transporters[name] = append(t.Transporters[name], transporters...)

	return nil
}

func (t *Task) RenderAll(ctx stdcontext.Context) error {
	return t.Render(ctx, t.AllocatedTemplateNameList(true)...)
}

func (t *Task) Render(ctx stdcontext.Context, names ...string) error {
	var buf bytes.Buffer
	var err error

	for _, name := range names {
		if !t.AllocatedTemplateExisted(name) {
			return fmt.Errorf("perform render failed as missing template `%s` in allocated", name)
		}
		err = t.tmpl.ExecuteTemplate(&buf, name, t.PerformCtx)
		if err != nil {
			return fmt.Errorf("perform render for template `%s` failed: %v", name, err)
		}

		err = t.PerformCtx.Setter(filepath.Join(name, t.RenderedEndTag), buf.String())
		if err != nil {
			return fmt.Errorf("perform render for template `%s` try to add rendered to Context failed: %v", name, err)
		}
		buf.Reset()
	}

	return nil
}

func (t *Task) Transport(ctx stdcontext.Context, groups ...string) error {
	if len(groups) == 0 {
		if len(t.Transporters) == 0 {
			return nil
		} else {
			for k := range t.Transporters {
				groups = append(groups, k)
			}
		}
	}

	sort.Strings(groups)

	var err error

	for _, group := range groups {
		trans := t.Transporters[group]
		if trans != nil {
			for idx, tran := range trans {
				err = tran.Transport(ctx, t.PerformCtx)
				if err != nil {
					return fmt.Errorf("transport failed for transporter index[%d] with in group `%s`: %v \n[%s]", idx, group, err, tran.Inspection())
				}
			}
		}
	}

	return nil
}

func (t *Task) Perform(ctx stdcontext.Context) error {

	err := t.RenderAll(ctx)
	if err != nil {
		return err
	}

	return t.Transport(ctx)
}
