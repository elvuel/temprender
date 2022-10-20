package main

import (
	stdctx "context"

	"github.com/elvuel/temprender/context"
	"github.com/elvuel/temprender/task"
	"github.com/elvuel/temprender/transport/puppet"
)

func main() {

	const tmpl = "quick.tmpl"
	t := task.NewTask()

	t.PerformCtx, _ = context.NewDefaultContext()
	t.PerformCtx.S("user", "Quick")
	t.PerformCtx.S("greeter", "temprender")

	t.LoadTempatesFromMap(map[string]string{
		tmpl: `Hi {{.G "user"}},

Greeting from {{.G "greeter"}}!`,
	})

	transporter, _ := puppet.NewPuppeteer()
	t.RegisterTransporters("default", transporter)

	t.RenderAll(stdctx.TODO())

	t.Transport(stdctx.TODO(), "default")

}
