package base

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/henrylee2cn/opay"
	"github.com/jmoiron/sqlx"
)

// Order base model
type (
	BaseOrder struct {
		Id string `json:"id"`
		//second party's order id
		Id2 string `json:"id2,omitempty"`
		//asset id
		Aid string `json:"aid"`
		//second party's asset id
		Aid2 string `json:"aid2,omitempty"`
		//user id
		Uid string `json:"uid"`
		//second party's user id
		Uid2 string `json:"uid2,omitempty"`
		//order type
		Type uint8 `json:"type"`
		//the amount of change for the Uid-Aid account, balance of positive and negative representation
		Amount float64 `json:"amount"`
		//the amount of change for the Uid-Aid account, balance of positive and negative representation
		Amount2 float64 `json:"amount2,omitempty"`

		Summary      string    `json:"summary"`
		Details      []*Detail `json:"details"`
		detailsBytes []byte
		Status       int32 `json:"status"`
		CreatedAt    int64 `json:"created_at"`

		//the most recent status
		lastStatus int32

		//processing error
		err error

		lock sync.RWMutex
	}
	Detail struct {
		UpdatedAt int64  `json:"updated_at"`
		Status    int32  `json:"status"`
		Notes     string `json:"notes"`
		Ip        string `json:"ip"`
	}
)

var _ opay.IOrder = new(BaseOrder)

// Prepare order status before push opay.
func (this *BaseOrder) Prepare(status int32, notes string, ip string) *BaseOrder {
	this.lastStatus = this.Status
	this.Status = status
	return this
}

// Add order detail after prepared status.
func (this *BaseOrder) AddDetail(notes string, ip string) *BaseOrder {
	if len(notes) == 0 {
		notes = this.GetStatusText()
	}
	this.Details = append(this.Details, &Detail{
		UpdatedAt: time.Now().Unix(),
		Status:    this.Status,
		Notes:     notes,
		Ip:        ip,
	})
	return this
}

// Rollback order status and detail in memory after dealing failure.
func (this *BaseOrder) Rollback() *BaseOrder {
	count := len(this.Details)
	if this.Details[count-1].Status == this.Status {
		this.Details = this.Details[:count-1]
	}
	this.detailsBytes = nil
	this.Status = this.lastStatus
	return this
}

// Get details of the bytes format.
func (this *BaseOrder) DetailsBytes(status int32, notes string, ip string) []byte {
	if this.detailsBytes == nil {
		if this.Details == nil {
			this.Details = []*Detail{}
		}
		this.detailsBytes, _ = json.Marshal(this.Details)
	}

	return this.detailsBytes
}

// Get the most recent Action, the default value is UNSET==0.
func (this *BaseOrder) LastAction() opay.Action {
	return opay.Action(this.lastStatus)
}

// Get user's id.
func (this *BaseOrder) GetUid() string {
	return this.Uid
}

// Get the second party's user id.
func (this *BaseOrder) GetUid2() string {
	return this.Uid2
}

// Get asset id.
func (this *BaseOrder) GetAid() string {
	return this.Aid
}

// Get the second party's asset id. (for example, the currency exchange business)
func (this *BaseOrder) GetAid2() string {
	return this.Aid2
}

// Get the amount of change for the Uid-Aid account,
// balance of positive and negative representation.
func (this *BaseOrder) GetAmount() float64 {
	return this.Amount
}

// Get the amount of change for the Uid-Aid2 account,
// balance of positive and negative representation.
func (this *BaseOrder) GetAmount2() float64 {
	return this.Amount2
}

// Async execution, and mark pending.
func (this *BaseOrder) ToPend(tx *sqlx.Tx, values opay.Values) error {
	return errors.New("*BaseOrder does not implement opay.IOrder (missing ToPend method).")
}

// Async execution, and mark the doing.
func (this *BaseOrder) ToDo(tx *sqlx.Tx, values opay.Values) error {
	return errors.New("*BaseOrder does not implement opay.IOrder (missing ToDo method).")
}

// Async execution, and mark the successful.
func (this *BaseOrder) ToSucceed(tx *sqlx.Tx, values opay.Values) error {
	return errors.New("*BaseOrder does not implement opay.IOrder (missing ToSucceed method).")
}

// Async execution, and mark canceled.
func (this *BaseOrder) ToCancel(tx *sqlx.Tx, values opay.Values) error {
	return errors.New("*BaseOrder does not implement opay.IOrder (missing ToCancel method).")
}

// Async execution, and mark failure.
func (this *BaseOrder) ToFail(tx *sqlx.Tx, values opay.Values) error {
	return errors.New("*BaseOrder does not implement opay.IOrder (missing ToFail method).")
}

// Sync execution, and mark the successful.
func (this *BaseOrder) SyncDeal(tx *sqlx.Tx, values opay.Values) error {
	return errors.New("*BaseOrder does not implement opay.IOrder (missing SyncDeal method).")
}

// Get the order's id.
func (this *BaseOrder) GetId() string {
	return this.Id
}

// Get the second party's order id.
func (this *BaseOrder) GetId2() string {
	return this.Id2
}

// Get the order's type.
func (this *BaseOrder) GetType() uint8 {
	return this.Type
}

// Get status text.
func (this *BaseOrder) GetStatusText() string {
	return GetStatusText(this.Type, this.Status)
}

// Get the order's summary.
func (this *BaseOrder) GetSummary() string {
	return this.Summary
}

// Get the order processing record details.
func (this *BaseOrder) GetDetails() []*Detail {
	return this.Details
}

// Get the order's status.
func (this *BaseOrder) GetStatus() int32 {
	return this.Status
}

// Get the order's created time.
func (this *BaseOrder) GetCreatedAt() int64 {
	return this.CreatedAt
}

// Binding the order and it's related order.
func (this *BaseOrder) Bind(other *BaseOrder) {
	this.Id2, this.Uid2 = other.Id, other.Uid
	other.Id2, other.Uid2 = this.Id, this.Uid
}

// set order id, 32bytes(time23+type3+random6)
func (this *BaseOrder) setId() *BaseOrder {
	this.Id = CreateOrderid(this.Type)
	return this
}
