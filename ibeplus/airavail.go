package ibeplus

import (
	"encoding/xml"
	"fmt"

	"github.com/otwdev/ibepluslib/models"

	"github.com/otwdev/galaxylib"
)

const airAvailURL = "http://ibeplus.travelsky.com/ota/xml/AirAvail/RealTime"

type AirAvail struct {
	Order    *models.OrderInfo
	Quantity string
}

func NewAirAvail(order *models.OrderInfo) *AirAvail {
	return &AirAvail{
		Order: order,
	}
}

func (a *AirAvail) AirAvailRQ() *galaxylib.GalaxyError {
	data := a.bindRQData()
	ibe := NewIBE(airAvailURL, data.RQAvailOTA_AirAvailRQ)
	rsData, err := ibe.Reqeust() //ReqeustIBE(airAvailURL, data.RQAvailOTA_AirAvailRQ)
	if err != nil {
		return err
	}
	fmt.Println(string(rsData))

	rs := &RSAvailRoot{}
	rs.RSAvailOTA_AirAvailRS = &RSAvailOTA_AirAvailRS{}
	if err := xml.Unmarshal(rsData, rs.RSAvailOTA_AirAvailRS); err != nil {
		return galaxylib.DefaultGalaxyError.FromError(1, err) //utils.FromError(1, err)
	}

	if rs.RSAvailOTA_AirAvailRS.RSAvailErrors != nil && rs.RSAvailOTA_AirAvailRS.RSAvailErrors.RSAvailError != nil {
		return galaxylib.DefaultGalaxyError.FromText(1, rs.RSAvailOTA_AirAvailRS.RSAvailErrors.RSAvailError.AttrShortTextZH) //utils.NewError(1, rs.RSAvailOTA_AirAvailRS.RSAvailErrors.RSAvailError.AttrShortTextZH)
	}

	return nil
}

func (a *AirAvail) bindRQData() *RQAvailRoot {
	pnr := a.Order.PnrInofs[0]
	rq := &RQAvailRoot{}
	rq.RQAvailOTA_AirAvailRQ = &RQAvailOTA_AirAvailRQ{}
	rq.RQAvailOTA_AirAvailRQ.RQAvailPOS = &RQAvailPOS{}
	rq.RQAvailOTA_AirAvailRQ.RQAvailPOS.RQAvailSource = &RQAvailSource{
		AttrPseudoCityCode: pnr.OfficeNumber,
	}

	rq.RQAvailOTA_AirAvailRQ.RQAvailOriginDestinationInformation = &RQAvailOriginDestinationInformation{}
	rq.RQAvailOTA_AirAvailRQ.RQAvailOriginDestinationInformation.RQAvailOriginDestinationOptions = &RQAvailOriginDestinationOptions{}
	rq.RQAvailOTA_AirAvailRQ.RQAvailOriginDestinationInformation.RQAvailOriginDestinationOptions.RQAvailOriginDestinationOption = &RQAvailOriginDestinationOption{}
	//rq.RQAvailOTA_AirAvailRQ.RQAvailOriginDestinationInformation.RQAvailOriginDestinationOptions.RQAvailOriginDestinationOption

	for _, v := range pnr.FlightSegments {
		segment := &RQAvailFlightSegment{
			AttrFlightNumber:      v.FlyNo,
			AttrDepartureDateTime: v.DepartrueDateTime(),
			RQAvailDepartureAirport: &RQAvailDepartureAirport{
				AttrLocationCode: v.DepartCityCode,
			},
			RQAvailArrivalAirport: &RQAvailArrivalAirport{
				AttrLocationCode: v.ArriveCityCode,
			},
			RQAvailMarketingAirline: &RQAvailMarketingAirline{
				AttrCode: v.MarketingAirLine,
			},
		}
		rq.RQAvailOTA_AirAvailRQ.RQAvailOriginDestinationInformation.RQAvailOriginDestinationOptions.RQAvailOriginDestinationOption.RQAvailFlightSegment =
			append(rq.RQAvailOTA_AirAvailRQ.RQAvailOriginDestinationInformation.RQAvailOriginDestinationOptions.RQAvailOriginDestinationOption.RQAvailFlightSegment, segment)
	}
	return rq
}

/********************Response*******************************/

type RSAvailRoot struct {
	RSAvailOTA_AirAvailRS *RSAvailOTA_AirAvailRS `xml:" OTA_AirAvailRS,omitempty" json:"OTA_AirAvailRS,omitempty"`
}

type RSAvailError struct {
	AttrCode        string        `xml:" Code,attr"  json:",omitempty"`
	AttrShortText   string        `xml:" ShortText,attr"  json:",omitempty"`
	AttrShortTextZH string        `xml:" ShortTextZH,attr"  json:",omitempty"`
	AttrType        string        `xml:" Type,attr"  json:",omitempty"`
	RSAvailTrace    *RSAvailTrace `xml:" Trace,omitempty" json:"Trace,omitempty"`
}

type RSAvailErrors struct {
	RSAvailError *RSAvailError `xml:" Error,omitempty" json:"Error,omitempty"`
}

type RSAvailOTA_AirAvailRS struct {
	RSAvailErrors *RSAvailErrors `xml:" Errors,omitempty" json:"Errors,omitempty"`
}

type RSAvailTrace struct {
	AttrText string `xml:" Text,attr"  json:",omitempty"`
}

/*******************Request****************************/

type RQAvailArrivalAirport struct {
	AttrLocationCode string `xml:" LocationCode,attr"  json:",omitempty"`
}

type RQAvailRoot struct {
	RQAvailOTA_AirAvailRQ *RQAvailOTA_AirAvailRQ `xml:" OTA_AirAvailRQ,omitempty" json:"OTA_AirAvailRQ,omitempty"`
}

type RQAvailDepartureAirport struct {
	AttrLocationCode string `xml:" LocationCode,attr"  json:",omitempty"`
}

type RQAvailFlightSegment struct {
	AttrDepartureDateTime   string                   `xml:" DepartureDateTime,attr"  json:",omitempty"`
	AttrFlightNumber        string                   `xml:" FlightNumber,attr"  json:",omitempty"`
	RQAvailArrivalAirport   *RQAvailArrivalAirport   `xml:" ArrivalAirport,omitempty" json:"ArrivalAirport,omitempty"`
	RQAvailDepartureAirport *RQAvailDepartureAirport `xml:" DepartureAirport,omitempty" json:"DepartureAirport,omitempty"`
	RQAvailMarketingAirline *RQAvailMarketingAirline `xml:" MarketingAirline,omitempty" json:"MarketingAirline,omitempty"`
}

type RQAvailMarketingAirline struct {
	AttrCode string `xml:" Code,attr"  json:",omitempty"`
}

type RQAvailOTA_AirAvailRQ struct {
	RQAvailOriginDestinationInformation *RQAvailOriginDestinationInformation `xml:" OriginDestinationInformation,omitempty" json:"OriginDestinationInformation,omitempty"`
	RQAvailPOS                          *RQAvailPOS                          `xml:" POS,omitempty" json:"POS,omitempty"`
	XMLName                             xml.Name                             `xml:" OTA_AirAvailRQ"`
}

type RQAvailOriginDestinationInformation struct {
	RQAvailOriginDestinationOptions *RQAvailOriginDestinationOptions `xml:" OriginDestinationOptions,omitempty" json:"OriginDestinationOptions,omitempty"`
}

type RQAvailOriginDestinationOption struct {
	RQAvailFlightSegment []*RQAvailFlightSegment `xml:" FlightSegment,omitempty" json:"FlightSegment,omitempty"`
}

type RQAvailOriginDestinationOptions struct {
	RQAvailOriginDestinationOption *RQAvailOriginDestinationOption `xml:" OriginDestinationOption,omitempty" json:"OriginDestinationOption,omitempty"`
}

type RQAvailPOS struct {
	RQAvailSource *RQAvailSource `xml:" Source,omitempty" json:"Source,omitempty"`
}

type RQAvailSource struct {
	AttrPseudoCityCode string `xml:" PseudoCityCode,attr"  json:",omitempty"`
}
