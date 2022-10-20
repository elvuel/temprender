package context

import (
	"encoding/json"
	"fmt"
)

var (
	_ Context = (*QuickContext)(nil)
)

type QuickContext struct {
	hostContext string
	Items       map[string]interface{} `json:"items,omitempty"`
}

func NewQuickContext(hostctx string) *QuickContext {
	return &QuickContext{hostContext: hostctx, Items: make(map[string]interface{})}
}

// Kind implements Context interface
func (ctx *QuickContext) Kind() string {
	return ctx.hostContext
}

// PanicableGetter implements Context AttrAccessor interface
func (ctx *QuickContext) PanicableGetter() bool {
	return true
}

// HasKey implements Context AttrAccessor interface
func (ctx *QuickContext) HasKey(k string) bool {
	_, ok := ctx.Items[k]
	return ok
}

// Setter implements Context AttrAccessor interface
func (ctx *QuickContext) Setter(k string, val interface{}) error {
	ctx.Items[k] = val
	return nil
}

// S implements Context AttrAccessor interface also a short alias form Setter
func (ctx *QuickContext) S(k string, val interface{}) error {
	return ctx.Setter(k, val)
}

// Getter implements Context AttrAccessor interface
func (ctx *QuickContext) Getter(k string) interface{} {
	val, ok := ctx.Items[k]

	if !ok && ctx.PanicableGetter() {
		panic(fmt.Sprintf("the `%s` is missing from the context [%s]", k, ctx.hostContext))
	}

	return val
}

// G implements Context AttrAccessor interface also a short alias form Getter
func (ctx *QuickContext) G(k string) interface{} {
	return ctx.Getter(k)
}

// Marshal implements Context marshaler interface{}
func (ctx *QuickContext) Marshal(interface{}) ([]byte, error) {
	return json.Marshal(ctx.Items)
}

// Unmarshal implements Context unmarshaler interface{}
func (ctx *QuickContext) Unmarshal(data []byte, _ interface{}) error {
	return json.Unmarshal(data, &ctx.Items)
}

// for puppet
func (ctx *QuickContext) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &ctx.Items)
}
