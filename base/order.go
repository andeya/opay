package base

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/henrylee2cn/opay"
	"github.com/jmoiron/sqlx"
)

type (
	// Order base model
	BaseOrder struct {
		Id        string `json:"id" db:"id"`
		LinkIdAid string `json:"link_id_aid" db:"link_id_aid"`
		Aid       string `json:"aid" db:"aid"`   //asset id
		Uid       string `json:"uid" db:"uid"`   //user id
		Type      string `json:"type" db:"type"` //order type
		//the amount of change for the Uid-Aid account, balance of positive and negative representation
		Amount        float64 `json:"amount" db:"amount"`
		Summary       string  `json:"summary" db:"summary"`
		Details       Details `json:"details" db:"details"`
		detailsString string
		preStatus     int64 //the previous status
		Status        int64 `json:"status" db:"status"` //the target status
		CreatedAt     int64 `json:"created_at" db:"created_at"`

		meta *opay.Meta
		err  error //processing error
		lock sync.RWMutex
	}
	Details []*Detail
	Detail  struct {
		UpdatedAt int64  `json:"updated_at" db:"-"`
		Status    int64  `json:"status" db:"-"`
		Note      string `json:"note" db:"-"`
		Ip        string `json:"ip" db:"-"`
	}
)

var _ opay.IOrder = new(BaseOrder)

//note: if param note is empty, do not append detail;
//and if param id is empty, the BaseOrder is new one.
func NewBaseOrder(
	meta *opay.Meta,
	aid string,
	uid string,
	amount float64,
	summary string,
	targetStatus int64,
	ip string,
) (*BaseOrder, error) {
	if meta == nil {
		return nil, errors.New("Param meta can not be nil.")
	}
	_, ok := meta.Status(targetStatus)
	if !ok {
		return nil, errors.New("Target status is invalid.")
	}
	if len(aid) == 0 || len(aid) > 2 || strings.HasPrefix(aid, "0") {
		return nil, errors.New("wrong aid format.")
	}
	var o = &BaseOrder{
		Id:        createOrderid(aid),
		CreatedAt: time.Now().Unix(),
		Aid:       aid,
		Uid:       uid,
		Type:      meta.OrderType(),
		Amount:    amount,
		Summary:   summary,
		Status:    meta.UnsetCode(),
		Details:   []*Detail{},
		meta:      meta,
	}
	err := o.SetTarget(targetStatus, ip)
	if err != nil {
		return nil, err
	}
	return o, nil
}

// Specify the handler of dealing.
func (this *BaseOrder) GetMeta() *opay.Meta {
	return this.meta
}

// Specify the handler of dealing.
func (this *BaseOrder) SetMeta(meta *opay.Meta) error {
	if meta == nil {
		return errors.New("Param meta can not be nil.")
	}
	this.meta = meta
	return nil
}

func (this *BaseOrder) PreStatus() int64 {
	return this.preStatus
}

// Get the order's target status.
func (this *BaseOrder) TargetStatus() int64 {
	return this.Status
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
func (this *BaseOrder) Pend(tx *sqlx.Tx) error {
	return errors.New("*BaseOrder does not implement opay.IOrder (missing Pend method).")
}

// Async execution, and mark the doing.
func (this *BaseOrder) Do(tx *sqlx.Tx) error {
	return errors.New("*BaseOrder does not implement opay.IOrder (missing Do method).")
}

// Async execution, and mark the successful.
func (this *BaseOrder) Succeed(tx *sqlx.Tx) error {
	return errors.New("*BaseOrder does not implement opay.IOrder (missing Succeed method).")
}

// Async execution, and mark canceled.
func (this *BaseOrder) Cancel(tx *sqlx.Tx) error {
	return errors.New("*BaseOrder does not implement opay.IOrder (missing Cancel method).")
}

// Async execution, and mark failure.
func (this *BaseOrder) Fail(tx *sqlx.Tx) error {
	return errors.New("*BaseOrder does not implement opay.IOrder (missing Fail method).")
}

// Sync execution, and mark the successful.
func (this *BaseOrder) SyncDeal(tx *sqlx.Tx) error {
	return errors.New("*BaseOrder does not implement opay.IOrder (missing SyncDeal method).")
}

// Set the target Action.
func (this *BaseOrder) SetTarget(targetStatus int64, ip string) error {
	if this.Status == targetStatus {
		return errors.New("Target status and the current status is the same.")
	}
	this.preStatus, this.Status = this.Status, targetStatus

	if this.Details == nil {
		this.Details = []*Detail{}
	}
	this.Details = append(this.Details, &Detail{
		UpdatedAt: time.Now().Unix(),
		Status:    this.Status,
		Note:      this.meta.Note(this.Status),
		Ip:        ip,
	})
	return nil
}

// Binding the order and it's related order.
func (this *BaseOrder) Link(related *BaseOrder) {
	this.LinkIdAid, related.LinkIdAid = related.Id+"|"+related.Aid, this.Id+"|"+this.Aid
}

// Get the related order's 'id' and 'aid'.
func (this *BaseOrder) SplitLink() (id, aid string) {
	if len(this.LinkIdAid) == 0 {
		return "", this.Aid
	}
	a := strings.Split(this.LinkIdAid, "|")
	if len(a) != 2 {
		return "", ""
	}
	return a[0], a[1]
}

// Get details of the json string format.
func (this *BaseOrder) DetailsString() string {
	if len(this.detailsString) == 0 {
		s, _ := this.Details.Value()
		this.detailsString = s.(string)
	}
	return this.detailsString
}

// Rollback order status and detail in memory after dealing failure.
func (this *BaseOrder) Rollback() *BaseOrder {
	count := len(this.Details)
	if count > 0 && this.Details[count-1].Status == this.Status {
		this.Details = this.Details[:count-1]
	}
	this.detailsString = ""
	this.Status = this.preStatus
	return this
}

// Get the order's id.
func (this *BaseOrder) GetId() string {
	return this.Id
}

// Get the order's summary.
func (this *BaseOrder) GetSummary() string {
	return this.Summary
}

// Get the order processing record details.
func (this *BaseOrder) GetDetails() []*Detail {
	return this.Details
}

// Get the order's created time.
func (this *BaseOrder) GetCreatedAt() int64 {
	return this.CreatedAt
}

var (
	_ sql.Scanner   = &Details{}
	_ driver.Valuer = &Details{}
)

// Scan implements the sql Scanner interface.
func (this *Details) Scan(value interface{}) error {
	v, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Cannot convert 'details' type %T to type 'Details'.", value)
	}
	if len(v) == 0 {
		if this != nil {
			*this = Details{}
		}
		return nil
	}
	// debug
	// println(string(([]byte)(v)))
	return json.Unmarshal(v, this)
}

// Value implements the driver Valuer interface.
func (this *Details) Value() (driver.Value, error) {
	if this == nil {
		return "[]", nil
	}
	b, err := json.Marshal(this)
	// debug
	// println(string(([]byte)(b)))
	return *(*string)(unsafe.Pointer(&b)), err
}
