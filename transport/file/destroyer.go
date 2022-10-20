package file

import (
	stdcontext "context"
	"os"

	"github.com/elvuel/temprender/context"
	"github.com/elvuel/temprender/transport"
)

type Destroyer struct {
	Kind   string `json:"tr_trans_kind,omitempty"`
	Target string `json:"target,omitempty"`
}

func NewDestroyerRegister() (transport.Transporter, error) {
	return NewDestroyer()
}

func NewDestroyer() (*Destroyer, error) {
	return &Destroyer{Kind: FileDestroyerTransporterKind}, nil
}

func (trans *Destroyer) Transport(_ stdcontext.Context, ctx context.Context) error {
	return os.Remove(trans.Target)
}
