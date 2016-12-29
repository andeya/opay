package opay

type (
	// Handler is order processing interface, only function or structure types are allowed.
	Handler interface {
		ServeOpay(*Context) error
	}

	// HandlerFunc Order processing interface function
	HandlerFunc func(*Context) error
)

var _ Handler = HandlerFunc(nil)

// ServeOpay implements Handler interface.
func (hf HandlerFunc) ServeOpay(ctx *Context) error {
	return hf(ctx)
}
