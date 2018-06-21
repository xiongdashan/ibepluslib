package models

type OrderMode int

const (
	Normal  OrderMode = 0 //正常
	Change  OrderMode = 1 //改期
	UpCabin OrderMode = 2 //升舱
)

type OrderStatus int

const (
	ToBeConfirmed OrderStatus = 0 //待确认
	ToBeIssued    OrderStatus = 1 //已确认
	Issuing       OrderStatus = 2 //出票中
	Issued        OrderStatus = 3 //已出票

)

type OrderInfo struct {
	ObjOrderID            string
	ObjType               string
	UpStreamSettlePrice   float64
	DownStreamSettlePrice float64
	TotalPrice            float64
	UsaAddress            string
	TripType              string
	PnrInofs              []*PnrInfo
	ContactInfo           *ContactInfo
	ExtraData             string
	OrderID               string
	Mode                  int
	WaitTime              string
	Channel               string
	OfficeNumber          string
	Status                int
}
