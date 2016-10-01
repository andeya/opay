package opay

import (
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
)

type Opay struct {
	metas          map[string]*Meta
	queue          Queue    //request queue
	db             *sqlx.DB //global database operation instance
	*SettleFuncMap          //global map of SettleFunc
	*Floater
	metasLock sync.RWMutex
}

func NewOpay(db *sqlx.DB, queueCapacity int, numOfDecimalPlaces int) *Opay {
	opay := &Opay{
		SettleFuncMap: globalSettleFuncMap,
		db:            db,
		metas:         make(map[string]*Meta),
		Floater:       NewFloater(numOfDecimalPlaces),
	}
	opay.queue = newOrderChan(queueCapacity, opay)
	return opay
}

// 处理请求
func (opay *Opay) Do(req Request) *Response {
	return <-opay.queue.Push(req)
}

func (opay *Opay) DB() *sqlx.DB {
	return opay.db
}

// Opay start.
func (opay *Opay) Serve() {
	if err := opay.db.Ping(); err != nil {
		panic(err)
	}
	for {
		// 读出一条请求
		// 无限等待
		req := opay.queue.Pull()

		var err error

		// 获取相应资产类型的账户余额操作函数
		var (
			initiatorSettle   SettleFunc
			stakeholderSettle SettleFunc
		)

		initiatorSettle, err = opay.GetSettleFunc(req.Initiator.GetAid())
		if err != nil {
			// 指定的资产账户的操作接口不存在时返回
			req.setError(err)
			req.writeback()
			continue
		}
		if req.Stakeholder != nil {
			stakeholderSettle, err = opay.GetSettleFunc(req.Stakeholder.GetAid())
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
				req.Tx, err = opay.db.Beginx()
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

			err = req.Initiator.GetMeta().serve(&Context{
				initiatorSettle:   initiatorSettle,
				stakeholderSettle: stakeholderSettle,
				Request:           req,
				Response:          req.response,
				Floater:           opay.Floater,
			})
		}()
	}
}
