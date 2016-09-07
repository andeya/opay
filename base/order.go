package base

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"sync"
	"time"
	"unsafe"

	"github.com/henrylee2cn/opay"
	"github.com/jmoiron/sqlx"
)

type (
	// Order base model
	BaseOrder struct {
		Id     string `json:"id" db:"id"`
		LinkId string `json:"link_id" db:"link_id"`
		Aid    string `json:"aid" db:"aid"`   //asset id
		Uid    string `json:"uid" db:"uid"`   //user id
		Type   uint8  `json:"type" db:"type"` //order type
		//the amount of change for the Uid-Aid account, balance of positive and negative representation
		Amount       float64 `json:"amount" db:"amount"`
		Summary      string  `json:"summary" db:"summary"`
		Details      Details `json:"details" db:"details"`
		detailsBytes []byte
		lastStatus   int32 //the most recent status
		Status       int32 `json:"status" db:"status"` //the target status
		CreatedAt    int64 `json:"created_at" db:"created_at"`

		err  error //processing error
		lock sync.RWMutex
	}
	Details []*Detail
	Detail  struct {
		UpdatedAt int64  `json:"updated_at" db:"-"`
		Status    int32  `json:"status" db:"-"`
		Note      string `json:"note" db:"-"`
		Ip        string `json:"ip" db:"-"`
	}
)

var _ opay.IOrder = new(BaseOrder)

//note: if param note is empty, do not append detail;
//and if param id is empty, the BaseOrder is new one.
func NewBaseOrder(
	id string,
	aid string,
	uid string,
	typ uint8,
	amount float64,
	summary string,
	curDetail []*Detail,
	curStatus int32,
	targetStatus int32,
	note string,
	ip string,
) *BaseOrder {
	var o = new(BaseOrder)
	t := time.Now().Unix()
	if len(id) == 0 {
		o.SetNewId()
		o.CreatedAt = t
	}
	o.Aid = aid
	o.Uid = uid
	o.Type = typ
	o.Amount = amount
	o.Summary = summary
	o.lastStatus = curStatus
	o.Status = targetStatus
	if curDetail == nil {
		o.Details = []*Detail{}
	} else {
		o.Details = curDetail
	}
	if len(note) > 0 {
		targetDetail := &Detail{
			UpdatedAt: t,
			Status:    o.Status,
			Note:      note,
			Ip:        ip,
		}
		o.Details = append(o.Details, targetDetail)
	}
	return o
}

// Specify the handler of dealing.
func (this *BaseOrder) Operator() string {
	return OrderOperator(this.Type)
}

// Get the target Action.
func (this *BaseOrder) TargetAction() opay.Action {
	return OrderAction(this.Type, this.Status)
}

// Get the most recent Action, the default value is UNSET==0.
func (this *BaseOrder) LastAction() opay.Action {
	return OrderAction(this.Type, this.lastStatus)
}

// Get user's id.
func (this *BaseOrder) GetUid() string {
	return this.Uid
}

// Get asset id.
func (this *BaseOrder) GetAid() string {
	return this.Aid
}

// Get the amount of change for the Uid-Aid account,
// balance of positive and negative representation.
func (this *BaseOrder) GetAmount() float64 {
	return this.Amount
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

// set order id, 32bytes(time23+type3+random6)
func (this *BaseOrder) SetNewId() *BaseOrder {
	this.Id = CreateOrderid(this.Type)
	return this
}

// Binding the order and it's related order.
func (this *BaseOrder) Link(related *BaseOrder) {
	this.LinkId, related.LinkId = related.Id, this.Id
}

// Get details of the bytes format.
func (this *BaseOrder) DetailsBytes() []byte {
	if this.detailsBytes == nil {
		if this.Details == nil {
			this.Details = []*Detail{}
		}
		this.detailsBytes, _ = json.Marshal(this.Details)
	}

	return this.detailsBytes
}

// Rollback order status and detail in memory after dealing failure.
func (this *BaseOrder) Rollback() *BaseOrder {
	count := len(this.Details)
	if count > 0 && this.Details[count-1].Status == this.Status {
		this.Details = this.Details[:count-1]
	}
	this.detailsBytes = nil
	this.Status = this.lastStatus
	return this
}

// Get the order's id.
func (this *BaseOrder) GetId() string {
	return this.Id
}

// Get the order's type.
func (this *BaseOrder) GetType() uint8 {
	return this.Type
}

// Get status text.
func (this *BaseOrder) GetStatusText() string {
	return OrderStatusText(this.Type, this.Status)
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

var (
	_ sql.Scanner   = new(Details)
	_ driver.Valuer = new(Details)
)

// Scan implements the sql Scanner interface.
func (this *Details) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	s, ok := value.(string)
	if !ok {
		return errors.New("cannot convert 'details' to type 'Details'.")
	}
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return json.Unmarshal(*(*[]byte)(unsafe.Pointer(&h)), this)
}

// Value implements the driver Valuer interface.
func (this *Details) Value() (driver.Value, error) {
	b, err := json.Marshal(this)
	return *(*string)(unsafe.Pointer(&b)), err
}
