package opay

import (
	"github.com/jmoiron/sqlx"
)

type Engine struct {
	*AccList           //账户操作接口列表
	*ServeMux          //订单操作接口全局路由
	Queue              //订单队列接口
	db        *sqlx.DB //全局数据库操作对象
}

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
		// 读出一条订单
		// 无限等待
		iOrd := engine.pull()

		// 获取账户操作接口
		var (
			accounter     Accounter
			withAccounter Accounter
			err           error
		)

		accounter, err = engine.GetAccounter(iOrd.GetAid())
		if err != nil {
			// 指定的资产账户的操作接口不存在时返回
			iOrd.Writeback(err)
		}

		withAccounter, err = engine.GetAccounter(iOrd.GetWithAid())
		if err != nil {
			// 指定的资产账户的操作接口不存在时返回
			iOrd.Writeback(err)
		}

		// 通过路由执行订单处理
		go engine.Exec(iOrd, accounter, withAccounter, engine.db)
	}
}
