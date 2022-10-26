package main

import (
	"bytes"
	stdctx "context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/elvuel/temprender/context"
	"github.com/elvuel/temprender/task"
	"github.com/elvuel/temprender/textual"
	filetransport "github.com/elvuel/temprender/transport/file"
)

const (
	baseDir = "./sample"
)

var (
	t *task.Task
)

func init() {
	t = task.NewTask()
	t.PerformCtx, _ = context.NewDefaultContext()
}

func main() {
	creator()
	injector()
	destroyer()
}

func creator() {
	const tmpl = "creator.tmpl"
	const group = "creator"
	creatorFixture(tmpl, group)
	ctx := stdctx.TODO()
	t.Render(ctx, tmpl)

	t.Transport(ctx, group)
}

func injector() {
	const tmpl = "injector.txt"
	const group = "injector"
	injectorFixture(tmpl, group)
	ctx := stdctx.TODO()
	t.Render(ctx, tmpl)

	t.Transport(ctx, group)
}

func destroyer() {
	const tmpl = "will_be_delete.txt"
	const group = "destroyer"
	destroyerFixture(tmpl, group)

	ctx := stdctx.TODO()
	t.Transport(ctx, group)
}

func creatorFixture(tmpl, group string) {
	targetDir := filepath.Join(baseDir, "creator")
	os.RemoveAll(targetDir)
	os.MkdirAll(targetDir, 0755)

	t.PerformCtx.S("user", "Creator")
	t.PerformCtx.S("greeter", "temprender")

	t.AppendTemplate(tmpl, `Hi {{.G "user"}}, Greeting from {{.G "greeter"}}!`)

	tran, _ := filetransport.NewCreator()
	tran.Key = filepath.Join(tmpl, t.RenderedEndTag)
	tran.Target = filepath.Join(targetDir, "greeting.txt")

	upcase := func(r io.Reader) (io.Reader, error) {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		data := buf.String()

		return bytes.NewBufferString(strings.ToUpper(data)), nil
	}

	tran.PreFitlers = make([]func(io.Reader) (io.Reader, error), 0)
	tran.PreFitlers = append(tran.PreFitlers, upcase)

	t.RegisterTransporters(group, tran)

}

func injectorFixture(tmpl, group string) {
	targetDir := filepath.Join(baseDir, "injector")
	os.RemoveAll(targetDir)
	os.MkdirAll(targetDir, 0755)

	injectTarget := filepath.Join(targetDir, tmpl)
	ioutil.WriteFile(injectTarget, []byte(`
Injector

placehoder1

injection1 - 1
injection2 - 2
injection3 - 3

placehoder2
`), 0644)

	t.PerformCtx.S("user", "Injector")
	t.PerformCtx.S("greeter", "temprender")
	t.PerformCtx.S("injector-only", "yup")

	t.AppendTemplate(tmpl, `Hi {{.G "user"}}, Greeting from {{.G "greeter"}}! -- {{ .G "injector-only" }}`)

	tran, _ := filetransport.NewInjector()
	tran.Target = injectTarget
	tran.Injections = make([]*filetransport.InjectionPattern, 0)

	placeholderExpr := "placehoder.[^\\n]*"
	injectionExpr := "injection[\\d]"
	tran.Injections = append(tran.Injections, &filetransport.InjectionPattern{
		Key: filepath.Join(tmpl, t.RenderedEndTag),
		Substitute: &textual.Substitute{
			Expr:         &placeholderExpr,
			Global:       true,
			FmtSpecifier: "%s",
		},
	}, &filetransport.InjectionPattern{
		FillWhenKeyMissing: true,
		Fills:              "_._",
		Substitute: &textual.Substitute{
			Expr:         &injectionExpr,
			Global:       true,
			FmtSpecifier: "%s",
		},
	})

	upcase := func(r io.Reader) (io.Reader, error) {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		data := buf.String()

		return bytes.NewBufferString(strings.ToUpper(data)), nil
	}

	tran.PreFitlers = make([]func(io.Reader) (io.Reader, error), 0)
	tran.PreFitlers = append(tran.PreFitlers, upcase)

	t.RegisterTransporters(group, tran)
}

func destroyerFixture(tmpl, group string) {
	targetDir := filepath.Join(baseDir, "destroyer")
	os.RemoveAll(targetDir)
	os.MkdirAll(targetDir, 0755)

	destroyTarget := filepath.Join(targetDir, tmpl)
	ioutil.WriteFile(destroyTarget, []byte(`
Destroyer

placehoder1

injection1 - 1
injection2 - 2
injection3 - 3

placehoder2
`), 0644)

	tran, _ := filetransport.NewDestroyer()
	tran.Target = destroyTarget

	t.RegisterTransporters(group, tran)
}
