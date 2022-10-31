package puppet

import (
	stdcontext "context"
	"encoding/json"
	"io"
	"os"

	"github.com/elvuel/temprender/context"
	"github.com/elvuel/temprender/transport"
)

const (
	PuppeteerTransporterKind = "debug::puppeteer"
)

func init() {
	transport.RegisterTransporter(&transport.TransporterManifest{
		Kind:    PuppeteerTransporterKind,
		NewFunc: NewPuppeteerRegister,
	})
}

type Puppeteer struct {
	Kind   string    `json:"tr_trans_kind,omitempty"`
	Writer io.Writer `json:"-"`
}

func NewPuppeteerRegister() (transport.Transporter, error) {
	return NewPuppeteer()
}

func NewPuppeteer() (*Puppeteer, error) {
	return &Puppeteer{Kind: PuppeteerTransporterKind, Writer: os.Stdout}, nil
}

func (trans *Puppeteer) Transport(_ stdcontext.Context, ctx context.Context) error {
	if trans.Writer == nil {
		trans.Writer = os.Stdout
	}

	encoder := json.NewEncoder(trans.Writer)
	encoder.SetIndent("", "\t")
	return encoder.Encode(ctx)
}

func (trans *Puppeteer) Inspection() string {
	return "puppeteer"
}
