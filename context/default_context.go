package context

var (
	_ Context = (*DefaultContext)(nil)

	defaultContextManifest = &ContextManifest{
		Kind:    defaultContext,
		NewFunc: NewDefaultContextRegister,
	}
)

type DefaultContext struct {
	Name string `json:"tr_ctx_kind,omitempty"`
	*QuickContext
}

func NewDefaultContextRegister() (Context, error) {
	return NewDefaultContext()
}

func NewDefaultContext() (*DefaultContext, error) {
	return &DefaultContext{Name: defaultContext, QuickContext: NewQuickContext(defaultContext)}, nil
}
