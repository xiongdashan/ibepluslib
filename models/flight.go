package models

import (
	"fmt"
	"strings"
	"time"
)

//FlightSegment 行程信息
type FlightSegment struct {
	DepartCityCode   string
	ArriveCityCode   string
	FlyNo            string
	Cabin            string
	ChdCabin         string
	Price            float64 //成人票面价
	ChdPrice         float64 //儿童票面价
	ChdRateFee       float64 //儿童燃油
	RateFee          float64 //燃油
	AirportFee       float64 //机建
	ChdAirportFee    float64
	BaseCabin        string //基础舱位
	FlyDate          string
	ArrDate          string
	DepTime          string
	ArrTime          string
	TripType         string
	TripSeq          int
	RuleRefund       string
	RuleChange       string
	PNRID            string
	OrderID          string
	FlightSegmentID  string
	MarketingAirLine string
	//ArriveDateTime string `gorm:`
}

func (v *FlightSegment) ArriveDateTime() string {

	return v.timeLayout(v.ArrDate, v.ArrTime)
}

func (v *FlightSegment) timeLayout(date, t string) string {

	dateTime := fmt.Sprintf("%s %s", date, t)

	result, _ := time.Parse("2006-01-02 15:04", dateTime)

	return result.Format("2006-01-02 15:04:00")
}

func (v *FlightSegment) DepartrueDateTime() string {
	return v.timeLayout(v.FlyDate, v.ArrTime)
}

func (v *FlightSegment) AirlineCode() string {
	return v.FlyNo[2:]
}

func (v *FlightSegment) AirlineNum() string {
	return v.FlyNo[:2]
}

func (f *FlightSegment) DepDate02Jan06() string {
	return f.date02Jan06(f.FlyDate)
}

func (f *FlightSegment) DepDate02Jan2006() string {
	return f.date02Jan2006(f.FlyDate)
}

func (f *FlightSegment) ArrDate02Jan2006() string {
	return f.date02Jan2006(f.ArrDate)
}

func (f *FlightSegment) ArrDate02Jan06() string {
	return f.date02Jan06(f.ArrDate)
}

func (f *FlightSegment) ArrTime1504() string {
	return f.time1504(f.ArrTime)
}

func (f *FlightSegment) DepTime1504() string {
	return f.time1504(f.DepTime)
}

func (f *FlightSegment) time1504(inputT string) string {
	t, _ := time.Parse("15:04", inputT)
	retT := t.Format("1504")
	return retT
}

func (f *FlightSegment) date02Jan06(inputD string) string {
	d, _ := time.Parse("2006-01-02", inputD)
	revD := d.Format("02Jan06")
	return strings.ToUpper(revD)
}

func (f *FlightSegment) date02Jan2006(inputD string) string {
	d, _ := time.Parse("2006-01-02", inputD)
	revD := d.Format("02Jan2006")
	return strings.ToUpper(revD)
}
