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
	var maxRoutine = opay.queue.GetCap() / 5
	if maxRoutine == 0 {
		maxRoutine = 1
	}
	var src = make(chan struct{}, maxRoutine)
	for {
		// Gets an execute permission
		src <- struct{}{}

		// Read a request
		// Unlimited wait
		req := opay.queue.Pull()

		var err error

		// Gets the account balance operation function for the corresponding asset type.
		var (
			initiatorSettle   SettleFunc
			stakeholderSettle SettleFunc
		)

		initiatorSettle, err = opay.GetSettleFunc(req.Initiator.GetAid())
		if err != nil {
			// Returns if the operation interface of the specified asset account does not exist.
			req.setError(err)
			req.writeback()
			continue
		}
		if req.Stakeholder != nil {
			stakeholderSettle, err = opay.GetSettleFunc(req.Stakeholder.GetAid())
			if err != nil {
				// Returns if the operation interface of the specified asset account does not exist
				req.setError(err)
				req.writeback()
				continue
			}
		}

		// The order processing is performed by routing.
		go func() {
			var err error
			defer func() {
				r := recover()
				if r != nil {
					err = fmt.Errorf("opay panic: %v", r)
				}

				// Close the request, and mark the end of the request processing
				req.setError(err)
				req.writeback()
				// Frees an execute permission
				<-src
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
