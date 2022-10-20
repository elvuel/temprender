package file

import (
	"bytes"
	stdcontext "context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/elvuel/temprender/context"
	"github.com/elvuel/temprender/transport"
)

type Strategy string

const (
	StrategyOverwrite Strategy = "overrite"
	StrategySkip      Strategy = "skip"
	StrategyAlias     Strategy = "alias"
)

type Creator struct {
	Kind            string   `json:"tr_trans_kind,omitempty"`
	Target          string   `json:"target,omitempty"`
	Key             string   `json:"target_ctx_key,omitempty"`
	ExistedStrategy Strategy `json:"existed_strategy,omitempty"` // target: overwrite, skip, alias
	AliasExt        string   `json:"alias_ext,omitempty"`
}

func NewCreatorRegister() (transport.Transporter, error) {
	return NewCreator()
}

func NewCreator() (*Creator, error) {
	return &Creator{Kind: FileCreatorTransporterKind, AliasExt: FileCreatorAliasExt}, nil
}

func (trans *Creator) Transport(_ stdcontext.Context, ctx context.Context) error {
	if _, err := os.Stat(trans.Target); os.IsNotExist(err) {
		// not exist
		goto createFile
	} else {
		// exists
		switch trans.ExistedStrategy {
		case StrategyOverwrite, "":
			goto createFile
		case StrategySkip:
			// log.Println("skipped to write file for", "`"+trans.Key+"`", "with", trans.Target)
			return nil
		case StrategyAlias:
			trans.Target += string(trans.AliasExt)
			goto createFile
		}
	}

createFile:
	lpath := filepath.Dir(trans.Target)
	os.MkdirAll(lpath, 0644)

	f, err := os.Create(trans.Target)
	if err != nil {
		return fmt.Errorf("file creator failed to create new file %s: %v", trans.Target, err)
	}

	data := ctx.Getter(trans.Key)

	var buf io.Reader
	switch data.(type) {
	case io.Reader:
		buf = data.(io.Reader)
	case string:
		buf = bytes.NewBufferString(data.(string))
	case []byte:
		buf = bytes.NewBuffer(data.([]byte))
	case nil:
		buf = bytes.NewBufferString("")
	default:
		return fmt.Errorf(
			"value type for context key[%s] in transporter[%s] should be one of [io.Reader, string, []byte, nil]",
			trans.Key, FileCreatorTransporterKind,
		)
	}

	_, err = io.Copy(f, buf)
	if err != nil {
		return fmt.Errorf("file creator failed to copy to new file %s: %v", trans.Target, err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("file creator failed to write to new file %s: %v", trans.Target, err)
	}

	return nil
}
