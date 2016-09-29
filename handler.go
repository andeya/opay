package opay

type (
	// 订单处理接口
	// 只允许函数或结构体类型
	Handler interface {
		ServeOpay(*Context) error
	}

	// 订单处理接口函数
	HandlerFunc func(*Context) error
)

var _ Handler = HandlerFunc(nil)

func (hf HandlerFunc) ServeOpay(ctx *Context) error {
	return hf(ctx)
}
