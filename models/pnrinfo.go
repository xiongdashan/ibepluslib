package models

type PNRStatus int64

const (
	PNRBookingFailured PNRStatus = 1
	PNRIssuing         PNRStatus = 2
	PNRIssued          PNRStatus = 4
	PNRCanceled        PNRStatus = 8
)

type PnrInfo struct {
	PnrCode        string
	BigPnr         string
	ChdPolicyID    string
	CommissionRate float64
	PnrText        string
	PolicyID       string
	InsureType     string
	InsurePrice    float64
	FlightSegments []*FlightSegment
	TravelerInfos  []*Traveler
	OrderID        string
	PnrID          string
	ADTQuantity    int
	CHDQuantity    int
	INFQuantity    int
	OfficeNumber   string
	IATANumber     string
	TripType       string
	Currency       string
	SegmentPrice   []*FlightSegmentFarePrice
	Status         int64
	Price          []*PnrPrice
	AirPortCode    string
	AgencyCity     string
	//	PNRSeq           uint `gorm:"type:varchar(30);index:idx_pnr_info_pnr_seq"`
}

func (p *PnrInfo) PersonQuantity() (int, int) {
	adt := 0
	chd := 0
	for _, v := range p.TravelerInfos {
		if v.Type == string(Adult) {
			adt++
		}
		if v.Type == string(Child) {
			chd++
		}
	}
	return adt, chd
}
