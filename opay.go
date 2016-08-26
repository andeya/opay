package opay

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type (
	// 订单处理引擎
	Engine struct {
		*AccList           //账户操作接口列表
		*ServeMux          //订单操作接口全局路由
		Queue              //订单队列接口
		db        *sqlx.DB //全局数据库操作对象
	}
)

// 新建订单处理服务
func NewOpay(db *sqlx.DB, queueCapacity int) *Engine {
	return &Engine{
		AccList:  globalAccList,
		ServeMux: globalServeMux,
		Queue:    newOrderChan(queueCapacity),
		db:       db,
	}
}

// 启动订单处理服务
func (engine *Engine) Serve() {
	if err := engine.db.Ping(); err != nil {
		panic(err)
	}
	for {
		// 读出一条请求
		// 无限等待
		req := engine.Queue.pull()

		var err error

		// 检查处理行为Action是否合法
		if err = req.ValidateAction(); err != nil {
			req.writeback(err)
			continue
		}

		// 获取账户操作接口
		var (
			accounter     Accounter
			withAccounter Accounter
		)

		accounter, err = engine.GetAccounter(req.IOrder.GetAid())
		if err != nil {
			// 指定的资产账户的操作接口不存在时返回
			req.writeback(err)
			continue
		}

		withAccounter, err = engine.GetAccounter(req.IOrder.GetWithAid())
		if err != nil {
			// 指定的资产账户的操作接口不存在时返回
			req.writeback(err)
			continue
		}

		// 通过路由执行订单处理
		go func() {
			var err error
			defer func() {
				r := recover()
				if r != nil {
					err = fmt.Errorf("%v", r)
				}

				// 关闭请求，标记请求处理结束
				req.writeback(err)
			}()

			if req.Tx == nil {
				req.Tx, err = engine.db.Beginx()
				if err != nil {
					return
				}
				defer func() {
					if err != nil {
						req.Tx.Rollback()
					} else {
						req.Tx.Commit()
					}
				}()
			}

			err = engine.ServeMux.serve(&Context{
				account:     accounter,
				withAccount: withAccounter,
				Request:     req,
			})
		}()
	}
}
