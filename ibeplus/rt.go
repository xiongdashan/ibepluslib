package ibeplus

import (
	"encoding/xml"

	"github.com/otwdev/ibepluslib/models"

	"github.com/otwdev/galaxylib"
)

const rtURL = "http://ibeplus.travelsky.com/ota/xml/AirResRet"

type RT struct {
	PNR *models.PnrInfo
}

func NewRT(pnr *models.PnrInfo) *RT {
	return &RT{pnr}
}

func (r *RT) RTPNR() (rs *RTRSOTA_AirResRetRS, err *galaxylib.GalaxyError) {
	rq := &RTRQOTA_AirResRetRQ{}
	rq.AttrRetCreateTimeInd = "true"
	rq.RTRQPOS = &RTRQPOS{}
	rq.RTRQPOS.RTRQSource = &RTRQSource{}
	rq.RTRQPOS.RTRQSource.AttrPseudoCityCode = r.PNR.OfficeNumber
	rq.RTRQBookingReferenceID = &RTRQBookingReferenceID{r.PNR.PnrCode}

	var ret []byte

	ibe := NewIBE(rtURL, rq)

	ret, err = ibe.Reqeust()
	if err != nil {
		return nil, err
	}

	//var rs *RTRSOTA_AirResRetRS

	if er := xml.Unmarshal(ret, &rs); er != nil {
		err = galaxylib.DefaultGalaxyError.FromError(1, err)
		return
	}

	if rs.RTRSErrors != nil {
		err = galaxylib.DefaultGalaxyError.FromText(1, rs.RTRSErrors.RTRSError.AttrShortTextZH)
		return
	}
	return
}

type RTRSAirResRet struct {
	RTRSAirTraveler             []*RTRSAirTraveler             `xml:" AirTraveler,omitempty" json:"AirTraveler,omitempty"`
	RTRSBookingReferenceID      *RTRSBookingReferenceID        `xml:" BookingReferenceID,omitempty" json:"BookingReferenceID,omitempty"`
	RTRSContactInfo             []*RTRSContactInfo             `xml:" ContactInfo,omitempty" json:"ContactInfo,omitempty"`
	RTRSCreateTime              *RTRSCreateTime                `xml:" CreateTime,omitempty" json:"CreateTime,omitempty"`
	RTRSFN                      *RTRSFN                        `xml:" FN,omitempty" json:"FN,omitempty"`
	RTRSFP                      *RTRSFP                        `xml:" FP,omitempty" json:"FP,omitempty"`
	RTRSFlightSegments          *RTRSFlightSegments            `xml:" FlightSegments,omitempty" json:"FlightSegments,omitempty"`
	RTRSOtherServiceInformation []*RTRSOtherServiceInformation `xml:" OtherServiceInformation,omitempty" json:"OtherServiceInformation,omitempty"`
	RTRSOthers                  []*RTRSOthers                  `xml:" Others,omitempty" json:"Others,omitempty"`
	RTRSResponsibility          *RTRSResponsibility            `xml:" Responsibility,omitempty" json:"Responsibility,omitempty"`
	RTRSSpecialRemark           []*RTRSSpecialRemark           `xml:" SpecialRemark,omitempty" json:"SpecialRemark,omitempty"`
	RTRSSpecialServiceRequest   []*RTRSSpecialServiceRequest   `xml:" SpecialServiceRequest,omitempty" json:"SpecialServiceRequest,omitempty"`
	RTRSTicketItemInfo          []*RTRSTicketItemInfo          `xml:" TicketItemInfo,omitempty" json:"TicketItemInfo,omitempty"`
	RTRSTicketing               *RTRSTicketing                 `xml:" Ticketing,omitempty" json:"Ticketing,omitempty"`
}

type RTRSAirTraveler struct {
	AttrRPH                   string                     `xml:" RPH,attr"  json:",omitempty"`
	RTRSPassengerTypeQuantity *RTRSPassengerTypeQuantity `xml:" PassengerTypeQuantity,omitempty" json:"PassengerTypeQuantity,omitempty"`
	RTRSPersonName            *RTRSPersonName            `xml:" PersonName,omitempty" json:"PersonName,omitempty"`
}

type RTRSAirline struct {
	AttrCode string `xml:" Code,attr"  json:",omitempty"`
}

type RTRSArrivalAirport struct {
	AttrLocationCode string `xml:" LocationCode,attr"  json:",omitempty"`
	AttrTerminal     string `xml:" Terminal,attr"  json:",omitempty"`
}

type RTRSBookingClassAvail struct {
	AttrResBookDesigCode string `xml:" ResBookDesigCode,attr"  json:",omitempty"`
}

type RTRSBookingReferenceID struct {
	AttrID string `xml:" ID,attr"  json:",omitempty"`
}

type RTRSRoot struct {
	RTRSOTA_AirResRetRS *RTRSOTA_AirResRetRS `xml:" OTA_AirResRetRS,omitempty" json:"OTA_AirResRetRS,omitempty"`
}

type RTRSContactInfo struct {
	AttrContactCity string `xml:" ContactCity,attr"  json:",omitempty"`
	AttrContactInfo string `xml:" ContactInfo,attr"  json:",omitempty"`
	AttrRPH         string `xml:" RPH,attr"  json:",omitempty"`
}

type RTRSCreateTime struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type RTRSDepartureAirport struct {
	AttrLocationCode string `xml:" LocationCode,attr"  json:",omitempty"`
	AttrTerminal     string `xml:" Terminal,attr"  json:",omitempty"`
}

type RTRSError struct {
	AttrCode        string     `xml:" Code,attr"  json:",omitempty"`
	AttrShortText   string     `xml:" ShortText,attr"  json:",omitempty"`
	AttrShortTextZH string     `xml:" ShortTextZH,attr"  json:",omitempty"`
	AttrType        string     `xml:" Type,attr"  json:",omitempty"`
	RTRSTrace       *RTRSTrace `xml:" Trace,omitempty" json:"Trace,omitempty"`
}

type RTRSErrors struct {
	RTRSError *RTRSError `xml:" Error,omitempty" json:"Error,omitempty"`
}

type RTRSFN struct {
	AttrRPH  string `xml:" RPH,attr"  json:",omitempty"`
	AttrText string `xml:" Text,attr"  json:",omitempty"`
}

type RTRSFP struct {
	AttrCurrency string `xml:" Currency,attr"  json:",omitempty"`
	AttrIsInfant string `xml:" IsInfant,attr"  json:",omitempty"`
	AttrPayType  string `xml:" PayType,attr"  json:",omitempty"`
	AttrRPH      string `xml:" RPH,attr"  json:",omitempty"`
	AttrRemark   string `xml:" Remark,attr"  json:",omitempty"`
}

type RTRSFlightSegment struct {
	AttrArrivalDateTime   string                 `xml:" ArrivalDateTime,attr"  json:",omitempty"`
	AttrDepartureDateTime string                 `xml:" DepartureDateTime,attr"  json:",omitempty"`
	AttrFlightNumber      string                 `xml:" FlightNumber,attr"  json:",omitempty"`
	AttrIsChanged         string                 `xml:" IsChanged,attr"  json:",omitempty"`
	AttrNumberInParty     string                 `xml:" NumberInParty,attr"  json:",omitempty"`
	AttrRPH               string                 `xml:" RPH,attr"  json:",omitempty"`
	AttrSegmentType       string                 `xml:" SegmentType,attr"  json:",omitempty"`
	AttrStatus            string                 `xml:" Status,attr"  json:",omitempty"`
	AttrTicket            string                 `xml:" Ticket,attr"  json:",omitempty"`
	RTRSArrivalAirport    *RTRSArrivalAirport    `xml:" ArrivalAirport,omitempty" json:"ArrivalAirport,omitempty"`
	RTRSBookingClassAvail *RTRSBookingClassAvail `xml:" BookingClassAvail,omitempty" json:"BookingClassAvail,omitempty"`
	RTRSDepartureAirport  *RTRSDepartureAirport  `xml:" DepartureAirport,omitempty" json:"DepartureAirport,omitempty"`
	RTRSMarketingAirline  *RTRSMarketingAirline  `xml:" MarketingAirline,omitempty" json:"MarketingAirline,omitempty"`
}

type RTRSFlightSegments struct {
	RTRSFlightSegment []*RTRSFlightSegment `xml:" FlightSegment,omitempty" json:"FlightSegment,omitempty"`
}

type RTRSMarketingAirline struct {
	AttrCode string `xml:" Code,attr"  json:",omitempty"`
}

type RTRSNamePNR struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type RTRSOTA_AirResRetRS struct {
	RTRSAirResRet *RTRSAirResRet `xml:" AirResRet,omitempty" json:"AirResRet,omitempty"`
	RTRSErrors    *RTRSErrors    `xml:" Errors,omitempty" json:"Errors,omitempty"`
}

type RTRSOtherServiceInformation struct {
	AttrRPH  string    `xml:" RPH,attr"  json:",omitempty"`
	RTRSText *RTRSText `xml:" Text,omitempty" json:"Text,omitempty"`
}

type RTRSOthers struct {
	AttrRPH  string `xml:" RPH,attr"  json:",omitempty"`
	AttrText string `xml:" Text,attr"  json:",omitempty"`
}

type RTRSPassengerTypeQuantity struct {
	AttrCode string `xml:" Code,attr"  json:",omitempty"`
}

type RTRSPersonName struct {
	RTRSNamePNR *RTRSNamePNR `xml:" NamePNR,omitempty" json:"NamePNR,omitempty"`
	RTRSSurname *RTRSSurname `xml:" Surname,omitempty" json:"Surname,omitempty"`
}

type RTRSResponsibility struct {
	AttrOfficeCode string `xml:" OfficeCode,attr"  json:",omitempty"`
	AttrRPH        string `xml:" RPH,attr"  json:",omitempty"`
}

type RTRSSpecialRemark struct {
	AttrRPH  string    `xml:" RPH,attr"  json:",omitempty"`
	RTRSText *RTRSText `xml:" Text,omitempty" json:"Text,omitempty"`
}

type RTRSSpecialServiceRequest struct {
	AttrRPH               string                 `xml:" RPH,attr"  json:",omitempty"`
	AttrSSRCode           string                 `xml:" SSRCode,attr"  json:",omitempty"`
	AttrStatus            string                 `xml:" Status,attr"  json:",omitempty"`
	RTRSAirline           *RTRSAirline           `xml:" Airline,omitempty" json:"Airline,omitempty"`
	RTRSText              *RTRSText              `xml:" Text,omitempty" json:"Text,omitempty"`
	RTRSTravelerRefNumber *RTRSTravelerRefNumber `xml:" TravelerRefNumber,omitempty" json:"TravelerRefNumber,omitempty"`
}

type RTRSSurname struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type RTRSText struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type RTRSTicketItemInfo struct {
	AttrRPH               string                 `xml:" RPH,attr"  json:",omitempty"`
	AttrTicketNumber      string                 `xml:" TicketNumber,attr"  json:",omitempty"`
	RTRSTravelerRefNumber *RTRSTravelerRefNumber `xml:" TravelerRefNumber,omitempty" json:"TravelerRefNumber,omitempty"`
}

type RTRSTicketing struct {
	AttrIsIssued   string `xml:" IsIssued,attr"  json:",omitempty"`
	AttrIssuedType string `xml:" IssuedType,attr"  json:",omitempty"`
	AttrRPH        string `xml:" RPH,attr"  json:",omitempty"`
	AttrRemark     string `xml:" Remark,attr"  json:",omitempty"`
}

type RTRSTrace struct {
	AttrText string `xml:" Text,attr"  json:",omitempty"`
}

type RTRSTravelerRefNumber struct {
	AttrRPH string `xml:" RPH,attr"  json:",omitempty"`
}

/********************************************
 RT RQ
********************************************/
type RTRQBookingReferenceID struct {
	AttrID string `xml:" ID,attr"  json:",omitempty"`
}

type RTRQRoot struct {
	RTRQOTA_AirResRetRQ *RTRQOTA_AirResRetRQ `xml:" OTA_AirResRetRQ,omitempty" json:"OTA_AirResRetRQ,omitempty"`
}

type RTRQOTA_AirResRetRQ struct {
	AttrRetCreateTimeInd   string                  `xml:" RetCreateTimeInd,attr"  json:",omitempty"`
	RTRQBookingReferenceID *RTRQBookingReferenceID `xml:" BookingReferenceID,omitempty" json:"BookingReferenceID,omitempty"`
	RTRQPOS                *RTRQPOS                `xml:" POS,omitempty" json:"POS,omitempty"`
}

type RTRQPOS struct {
	RTRQSource *RTRQSource `xml:" Source,omitempty" json:"Source,omitempty"`
}

type RTRQSource struct {
	AttrPseudoCityCode string `xml:" PseudoCityCode,attr"  json:",omitempty"`
}
