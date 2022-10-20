package task

import (
	"encoding/json"

	"github.com/elvuel/temprender/transport"
)

type GroupedTransporters map[string]transport.Transporters

func (g GroupedTransporters) UnmarshalJSON(data []byte) error {
	var rawMap map[string]*json.RawMessage
	err := json.Unmarshal(data, &rawMap)
	if err != nil {
		return nil
	}

	for k, v := range rawMap {
		trans, err := transport.UnmarshalTransporters(*v)
		if err != nil {
			return err
		}
		g[k] = trans
	}
	return nil
}
