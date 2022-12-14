package file

import (
	"bytes"
	stdcontext "context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/elvuel/temprender/context"
	"github.com/elvuel/temprender/textual"
	"github.com/elvuel/temprender/transport"
)

type Injector struct {
	Label     string  `json:"label,omitempty"`
	Kind      string  `json:"kind,omitempty"`
	Target    string  `json:"target,omitempty"`
	Fills     *string `json:"fills,omitempty"`
	WithAlias bool    `json:"with_alias,omitempty"`
	AliasExt  string  `json:"alias_ext,omitempty"`

	Injections []*InjectionPattern                  `json:"injections,omitempty"`
	PreFitlers []func(io.Reader) (io.Reader, error) `json:"-"` // apply adjust data filters before write
}

func NewInjectorRegister() (transport.Transporter, error) {
	return NewInjector()
}

func NewInjector() (*Injector, error) {
	return &Injector{Kind: FileInjectorTransporterKind, AliasExt: FileInjectorAliasExt}, nil
}

func (trans *Injector) Transport(_ stdcontext.Context, ctx context.Context) error {
	var loadedData string
	if _, err := os.Stat(trans.Target); os.IsNotExist(err) {
		if trans.Fills == nil {
			return errors.New("file injector's init file template data missing")
		} else {
			loadedData = *trans.Fills
		}
	} else {
		data, err := ioutil.ReadFile(trans.Target)
		if err != nil {
			return fmt.Errorf("file injector failed to read file %s: %v", trans.Target, err)
		}
		loadedData = string(data)
	}

	for _, pattern := range trans.Injections {
		if ctx.Exist(pattern.Key) {
			data := ctx.G(pattern.Key)

			var buf *bytes.Buffer
			switch data.(type) {
			case io.Reader:
				b, _ := io.ReadAll(data.(io.Reader))
				buf = bytes.NewBuffer(b)
			case string:
				buf = bytes.NewBufferString(data.(string))
			case []byte:
				buf = bytes.NewBuffer(data.([]byte))
			case nil:
				buf = bytes.NewBufferString("")
			default:
				return fmt.Errorf(
					"value type for context key[%s] in transporter[%s] should be one of [io.Reader, string, []byte, nil]",
					pattern.Key, FileInjectorTransporterKind,
				)
			}
			loadedData = pattern.Sub(loadedData, buf.String())
		} else {
			if pattern.FillWhenKeyMissing {
				loadedData = pattern.Sub(loadedData, pattern.Fills)
			} else {
				return fmt.Errorf("%s replacer corresponding context key %s missing", FileCreatorTransporterKind, pattern.Key)
			}
		}
	}

	if trans.WithAlias {
		trans.Target += string(trans.AliasExt)
	}

	lpath := filepath.Dir(trans.Target)

	os.MkdirAll(lpath, 0755)

	var buf io.Reader
	buf = bytes.NewBufferString(loadedData)

	var err error
	for fidx, filter := range trans.PreFitlers {
		buf, err = filter(buf)
		if err != nil {
			return fmt.Errorf("apply filter [%d] failed: %v", fidx, err)
		}
	}

	data, err := ioutil.ReadAll(buf)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(trans.Target, data, 0644)
}

func (trans *Injector) Inspection() string {
	if trans.Label == "" {
		return fmt.Sprintf("kind: %s for target %s, with alias?:%v", trans.Kind, trans.Target, trans.WithAlias)
	}
	return trans.Label
}

type InjectionPattern struct {
	Key                string `json:"key,omitempty"`
	FillWhenKeyMissing bool   `json:"fill_when_key_missing,omitempty"`
	Fills              string `json:"fills,omitempty"`

	*textual.Substitute
}
