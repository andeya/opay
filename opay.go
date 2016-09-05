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
		queue     Queue    //订单队列接口
		db        *sqlx.DB //全局数据库操作对象
	}
)

// 新建订单处理服务
func NewOpay(db *sqlx.DB, queueCapacity int) *Engine {
	return &Engine{
		AccList:  globalAccList,
		ServeMux: globalServeMux,
		queue:    newOrderChan(queueCapacity),
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
		req := engine.queue.Pull()

		var err error

		// 检查处理行为Action是否合法
		if err = req.ValidateAction(); err != nil {
			req.setError(err)
			req.writeback()
			continue
		}

		// 获取账户操作接口
		var (
			initiatorIAccount   IAccount
			stakeholderIAccount IAccount
		)

		initiatorIAccount, err = engine.GetIAccount(req.Initiator.GetAid())
		if err != nil {
			// 指定的资产账户的操作接口不存在时返回
			req.setError(err)
			req.writeback()
			continue
		}
		if req.Stakeholder != nil {
			stakeholderIAccount, err = engine.GetIAccount(req.Stakeholder.GetAid())
			if err != nil {
				// 指定的资产账户的操作接口不存在时返回
				req.setError(err)
				req.writeback()
				continue
			}
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
				req.setError(err)
				req.writeback()
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
				account:     initiatorIAccount,
				withAccount: stakeholderIAccount,
				Request:     req,
			})
		}()
	}
}

// 推送请求到引擎进行处理
func (engine *Engine) Push(req Request) (respChan <-chan Response, err error) {
	return engine.queue.Push(req)
}
