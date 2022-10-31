package transport

import (
	stdcontext "context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/elvuel/temprender/context"
)

const (
	KindTagName = "tr_trans_kind"
)

var (
	transporters map[string]NewTransporterFunc
)

type Transporter interface {
	Transport(stdcontext.Context, context.Context) error
	Inspection() string
}

type Transporters []Transporter

type NewTransporterFunc func() (Transporter, error)

type TransporterManifest struct {
	Kind    string
	NewFunc NewTransporterFunc
}

func init() {
	transporters = make(map[string]NewTransporterFunc)
}

// RegisterTransporter registers a transporter manifest
func RegisterTransporter(manifest *TransporterManifest) {
	transporters[manifest.Kind] = manifest.NewFunc
}

// NewTransporter returns a transporter
func NewTransporter(kind string) (Transporter, error) {
	funk, ok := transporters[kind]

	if !ok {
		return nil, fmt.Errorf("context manifests missing kind %s", kind)
	}

	return funk()
}

func UnmarshalTransporter(data []byte) (Transporter, error) {
	var rawMap map[string]*json.RawMessage
	err := json.Unmarshal(data, &rawMap)
	if err != nil {
		return nil, err
	}

	var tran Transporter

	for key, val := range rawMap {
		if key == KindTagName && val != nil {
			kval, _ := strconv.Unquote(string(*val))

			tran, err = NewTransporter(kval)
			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(data, tran)
			if err != nil {
				return nil, err
			}

			return tran, nil
		}
	}

	if tran != nil {
		return tran, nil
	}

	return nil, errors.New("invalid unmarshal data to unknown Transporter")
}

func UnmarshalTransporters(data []byte) (Transporters, error) {
	var slicer []*json.RawMessage
	err := json.Unmarshal(data, &slicer)
	if err != nil {
		return nil, err
	}

	trans := make(Transporters, 0)

	for _, item := range slicer {
		tran, err := UnmarshalTransporter(*item)
		if err != nil {
			return nil, err
		}
		trans = append(trans, tran)
	}

	return trans, nil
}
