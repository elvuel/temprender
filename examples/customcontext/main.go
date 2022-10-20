package main

import (
	stdctx "context"
	"time"

	"github.com/elvuel/temprender/context"
	"github.com/elvuel/temprender/task"
	"github.com/elvuel/temprender/transport/puppet"
)

const (
	customContextKind = "custom::context::example"
)

type CustomContext struct {
	CtxKind string `json:"tr_ctx_kind,omitempty"`
	User    string `json:"user,omitempty"`
	Greeter string `json:"greeter,omitempty"`
	*context.QuickContext
}

func NewCustomContextRegister() (context.Context, error) {
	return NewCustomContext()
}

func NewCustomContext() (*CustomContext, error) {
	return &CustomContext{CtxKind: customContextKind, QuickContext: context.NewQuickContext(customContextKind)}, nil
}

func init() {
	context.RegisterContext(&context.ContextManifest{
		Kind:    customContextKind,
		NewFunc: NewCustomContextRegister,
	})
}

func main() {
	const tmpl = "custom.tmpl"
	t := task.NewTask()

	customCtx, _ := NewCustomContext()
	customCtx.User = "Custom"
	customCtx.Greeter = "temprender"
	customCtx.S("timestamp", time.Now().Format(time.RFC3339))

	t.PerformCtx = customCtx

	t.LoadTempatesFromMap(map[string]string{
		tmpl: `Hi {{ .User }},

Greeting from {{ .Greeter }}!(with custom context {{ quote .CtxKind }})

Published at: {{ .G "timestamp" }}
`,
	})

	transporter, _ := puppet.NewPuppeteer()
	t.RegisterTransporters("default", transporter)

	t.RenderAll(stdctx.TODO())

	t.Transport(stdctx.TODO(), "default")
}
