package context

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

const (
	defaultContext = "default"

	KindTagName = "tr_ctx_kind"
)

var (
	contexts map[string]NewContextFunc
)

type AttrOperator interface {
	Setter(string, interface{}) error
	Getter(string) interface{}

	HasKey(string) bool

	S(string, interface{}) error
	G(string) interface{}

	PanicableGetter() bool
}

type Marshaler interface {
	Marshal(interface{}) ([]byte, error)
}

type Unmarshaler interface {
	Unmarshal([]byte, interface{}) error
}

// Context represents temprender data context
type Context interface {
	Kind() string

	Marshaler
	Unmarshaler

	AttrOperator
}

// NewContextFunc is type of Context initialization function
type NewContextFunc func() (Context, error)

type ContextManifest struct {
	Kind    string
	NewFunc NewContextFunc
}

func init() {
	contexts = make(map[string]NewContextFunc)

	RegisterContext(defaultContextManifest)
}

// RegisterContext registers a kind of Context initialization function
func RegisterContext(manifest *ContextManifest) {
	contexts[manifest.Kind] = manifest.NewFunc
}

// NewContext returns a new kind of Context
func NewContext(kind string) (Context, error) {
	funk, ok := contexts[kind]

	if !ok {
		return nil, fmt.Errorf("context manifests missing kind %s", kind)
	}

	return funk()
}

func Unmarshal(data []byte) (Context, error) {
	var rawMap map[string]*json.RawMessage
	err := json.Unmarshal(data, &rawMap)
	if err != nil {
		return nil, err
	}

	var ctx Context

	for key, val := range rawMap {
		if key == KindTagName && val != nil {
			kval, _ := strconv.Unquote(string(*val))

			ctx, err = NewContext(kval)
			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(data, ctx)
			if err != nil {
				return nil, err
			}

			return ctx, nil
		}
	}

	if ctx != nil {
		return ctx, nil
	}

	return nil, errors.New("invalid unmarshal data to unknown Context")
}
