package context

import (
	"encoding/json"
	"fmt"
	"strings"
)

var (
	_ Context = (*QuickContext)(nil)
)

type QuickContext struct {
	HostContext     string                 `json:"temprender_quickcontext_host_context_should_never_conflicts"`
	GetterPanicable bool                   `json:"temprender_quickcontext_getter_panicable_should_never_conflicts"`
	Items           map[string]interface{} `json:"temprender_quickcontext_items_should_never_conflicts"`
}

func NewQuickContext(hostctx string) *QuickContext {
	return &QuickContext{HostContext: hostctx, Items: make(map[string]interface{})}
}

// Kind implements Context interface returns context kind
func (ctx *QuickContext) Kind(args ...string) string {
	if len(args) > 0 {
		ctx.HostContext = strings.Join(args, "::")
	}
	return ctx.HostContext
}

// Clearn implements Context interface clear context
func (ctx *QuickContext) Clear() error {
	for k := range ctx.Items {
		delete(ctx.Items, k)
	}
	return nil
}

// PanicableGetter implements Context AttrAccessor interface
func (ctx *QuickContext) PanicableGetter() bool {
	return ctx.GetterPanicable
}

// Exist implements Context AttrAccessor interface checks key existed or not
func (ctx *QuickContext) Exist(key string) bool {
	_, ok := ctx.Items[key]
	return ok
}

// Delete implements Context AttrAccessor interface deletes key
func (ctx *QuickContext) Delete(k string) error {
	delete(ctx.Items, k)
	return nil
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
		panic(fmt.Sprintf("the `%s` is missing from the context [%s]", k, ctx.HostContext))
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
