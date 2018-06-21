package ibeplus

import (
	"encoding/xml"
	"ibepluslib/models"

	"github.com/otwdev/galaxylib"

	glinq "github.com/ahmetb/go-linq"
	valid "github.com/asaskevich/govalidator"
)

type AirFarePriceRequest struct {
	Order *models.OrderInfo
}

func NewFarePriceRequest(order *models.OrderInfo) *AirFarePriceRequest {
	return &AirFarePriceRequest{order}
}

//http://agibe.travelsky.com/ota/xml/AirFarePricePolicy/I

func (rq *AirFarePriceRequest) ValidFarePrice() (retErr *galaxylib.GalaxyError, combinNum int64) {

	if len(rq.Order.PnrInofs) > 1 {
		return galaxylib.DefaultGalaxyError.FromText(1, "只支持一个PNR"), 0
	}

	pnr := rq.Order.PnrInofs[0]

	fare := &FareRoot{}

	fare.FareTSK_AirfarePrice = &FareTSK_AirfarePrice{}
	fare.FareTSK_AirfarePrice.FareRequest = &FareRequest{}
	fare.FareTSK_AirfarePrice.FareRequest.FareSITA_AirfarePriceRQ = &FareSITA_AirfarePriceRQ{}

	fare.FareTSK_AirfarePrice.FareRequest.FareSITA_AirfarePriceRQ.FareOTA_AirPriceRQ = &FareOTA_AirPriceRQ{}

	fare.FareTSK_AirfarePrice.FareRequest.FareSITA_AirfarePriceRQ.FareOTA_AirPriceRQ.FarePOS = &FarePOS{}

	//Office 号
	// fare.FareTSK_AirfarePrice.FareRequest.FareSITA_AirfarePriceRQ.FareOTA_AirPriceRQ.FarePOS.FareSource =
	// 	append(fare.FareTSK_AirfarePrice.FareRequest.FareSITA_AirfarePriceRQ.FareOTA_AirPriceRQ.FarePOS.FareSource, &FareSource{
	// 		AttrPseudoCityCode: order.OfficeNumber,
	// 	})
	fare.FareTSK_AirfarePrice.FareRequest.FareSITA_AirfarePriceRQ.FareOTA_AirPriceRQ.FarePOS.FareSource =
		append(fare.FareTSK_AirfarePrice.FareRequest.FareSITA_AirfarePriceRQ.FareOTA_AirPriceRQ.FarePOS.FareSource, &FareSource{
			FareRequestorID: &FareRequestorID{
				AttrID:         pnr.IATANumber,
				AttrType:       "13",
				AttrID_Context: "IATA_Number",
			},
			AttrPseudoCityCode: pnr.OfficeNumber, //rq.Order.OfficeNumber,
			AttrAirportCode:    pnr.AirPortCode,  //"PEK",
		})

	// 航段
	fare.FareTSK_AirfarePrice.FareRequest.FareSITA_AirfarePriceRQ.FareOTA_AirPriceRQ.FareAirItinerary = &FareAirItinerary{}

	fare.FareTSK_AirfarePrice.FareRequest.FareSITA_AirfarePriceRQ.FareOTA_AirPriceRQ.FareAirItinerary.FareOriginDestinationOptions =
		adjustSegments(pnr) //a&FareOriginDestinationOptions{}

	//乘客
	fare.FareTSK_AirfarePrice.FareRequest.FareSITA_AirfarePriceRQ.FareOTA_AirPriceRQ.FareTravelerInfoSummary = &FareTravelerInfoSummary{}

	traveler := FareAirTravelerAvail{}

	//成人
	traveler = analysisTraveler()(traveler, models.Adult, pnr.ADTQuantity)

	//儿童
	traveler = analysisTraveler()(traveler, models.Child, pnr.CHDQuantity)

	// //婴儿
	// traveler = analysisTraveler()(traveler, air.Infant, pnr.INFQuantity)

	fare.FareTSK_AirfarePrice.FareRequest.FareSITA_AirfarePriceRQ.
		FareOTA_AirPriceRQ.FareTravelerInfoSummary.FareAirTravelerAvail = &traveler

	//价格信息
	fare.FareTSK_AirfarePrice.FareRequest.FareSITA_AirfarePriceRQ.
		FareOTA_AirPriceRQ.FareTravelerInfoSummary.FarePriceRequestInformation = &FarePriceRequestInformation{
		AttrCurrencyCode:  "CNY",
		AttrPricingSource: "Both",
		AttrReprice:       "false",
	}

	//额外信息
	fare.FareTSK_AirfarePrice.FareRequest.FareSITA_AirfarePriceRQ.FareAdditionalPriceRQData = &FareAdditionalPriceRQData{}
	fare.FareTSK_AirfarePrice.FareRequest.FareSITA_AirfarePriceRQ.FareAdditionalPriceRQData.
		AttrMaxResponses = "50"

	// fare.FareTSK_AirfarePrice.FareRequest.FareSITA_AirfarePriceRQ.FareAdditionalPriceRQData.
	// 	AttrReturnAllEndos = "true"

	fare.FareTSK_AirfarePrice.FareRequest.FareSITA_AirfarePriceRQ.FareAdditionalPriceRQData.
		FareTicketingCarrier = &FareTicketingCarrier{
		AttrCode: pnr.FlightSegments[0].MarketingAirLine,
	}

	fare.FareTSK_AirfarePrice.FareRequest.FareSITA_AirfarePriceRQ.FareOTA_AirPriceRQ.FareDistributor = &FareDistributor{
		Text: "运营二部/同业线下",
	}
	fare.FareTSK_AirfarePrice.FareRequest.FareSITA_AirfarePriceRQ.FareOTA_AirPriceRQ.FareGroupType = &FareGroupType{
		Text: "0",
	}

	rs := NewFarePriceResponse(pnr, fare)

	retErr, combinNum = rs.AirFarePriceRS() //AirFarePriceRS(fare)
	return
}

func analysisTraveler() func(traveler FareAirTravelerAvail, personType models.PersonType, quantity int) FareAirTravelerAvail {
	return func(traveler FareAirTravelerAvail, personType models.PersonType, quantity int) FareAirTravelerAvail {
		if quantity <= 0 {
			return traveler
		}
		f := &FarePassengerTypeQuantity{
			AttrCode:     string(personType),
			AttrQuantity: valid.ToString(quantity),
		}

		traveler.FarePassengerTypeQuantity = append(traveler.FarePassengerTypeQuantity, f)
		return traveler
	}
}

func adjustSegments(pnr *models.PnrInfo) *FareOriginDestinationOptions {

	revSegm := &FareOriginDestinationOptions{}
	//revSegm.FareOriginDestinationOption = &FareOriginDestinationOption{}

	var odOptions []*FareOriginDestinationOption

	glinq.From(pnr.FlightSegments).OrderByT(func(f *models.FlightSegment) int {
		return f.TripSeq
	}).ForEachT(func(v *models.FlightSegment) {
		fareSegm := &FareFlightSegment{
			AttrArrivalDateTime:   v.ArriveDateTime(),    //fmt.Sprintf("%sT%s", v.ArrDate, v.ArrTime),
			AttrDepartureDateTime: v.DepartrueDateTime(), //fmt.Sprintf("%sT%s", v.FlyDate, v.DepTime),
			AttrFlightNumber:      v.FlyNo,
		}
		fareSegm.FareDepartureAirport = &FareDepartureAirport{}
		fareSegm.FareArrivalAirport = &FareArrivalAirport{}
		fareSegm.FareArrivalAirport.AttrLocationCode = v.ArriveCityCode
		fareSegm.FareDepartureAirport.AttrLocationCode = v.DepartCityCode
		fareSegm.AttrResBookDesigCode = v.Cabin
		fareSegm.FareMarketingAirline = &FareMarketingAirline{
			AttrCode: v.MarketingAirLine,
		}
		// fareSegm.AttrRPH = valid.ToString(v.TripSeq) //strconv.FormatInt(intv.TripSeq, 10)
		opt := &FareOriginDestinationOption{}
		opt.FareFlightSegment = fareSegm

		odOptions = append(odOptions, opt)

	})
	revSegm.FareOriginDestinationOption = odOptions
	return revSegm
}

type FareAdditionalPriceRQData struct {
	// AttrETicketIndicator string                `xml:" ETicketIndicator,attr"  json:",omitempty"`
	AttrMaxResponses string `xml:" MaxResponses,attr"  json:",omitempty"`
	// AttrTaxBreakdownInd  string                `xml:" TaxBreakdownInd,attr"  json:",omitempty"`
	// AttrTaxSummaryInd    string                `xml:" TaxSummaryInd,attr"  json:",omitempty"`
	FareTicketingCarrier *FareTicketingCarrier `xml:" TicketingCarrier,omitempty" json:"TicketingCarrier,omitempty"`
}

type FareAirItinerary struct {
	FareOriginDestinationOptions *FareOriginDestinationOptions `xml:" OriginDestinationOptions,omitempty" json:"OriginDestinationOptions,omitempty"`
}

type FareAirTravelerAvail struct {
	FarePassengerTypeQuantity []*FarePassengerTypeQuantity `xml:" PassengerTypeQuantity,omitempty" json:"PassengerTypeQuantity,omitempty"`
}

type FareArrivalAirport struct {
	AttrCodeContext  string `xml:" CodeContext,attr"  json:",omitempty"`
	AttrLocationCode string `xml:" LocationCode,attr"  json:",omitempty"`
}

type FareRoot struct {
	FareTSK_AirfarePrice *FareTSK_AirfarePrice `xml:" TSK_AirfarePrice,omitempty" json:"TSK_AirfarePrice,omitempty"`
}

type FareDepartureAirport struct {
	AttrCodeContext  string `xml:" CodeContext,attr"  json:",omitempty"`
	AttrLocationCode string `xml:" LocationCode,attr"  json:",omitempty"`
}

type FareDistributor struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type FareEquipment struct {
	AttrAirEquipType string `xml:" AirEquipType,attr"  json:",omitempty"`
}

type FareFlightSegment struct {
	AttrArrivalDateTime   string `xml:" ArrivalDateTime,attr"  json:",omitempty"`
	AttrDepartureDateTime string `xml:" DepartureDateTime,attr"  json:",omitempty"`
	AttrFlightNumber      string `xml:" FlightNumber,attr"  json:",omitempty"`
	AttrResBookDesigCode  string `xml:" ResBookDesigCode,attr"  json:",omitempty"`
	// AttrStatus            string                `xml:" Status,attr"  json:",omitempty"`
	FareArrivalAirport   *FareArrivalAirport   `xml:" ArrivalAirport,omitempty" json:"ArrivalAirport,omitempty"`
	FareDepartureAirport *FareDepartureAirport `xml:" DepartureAirport,omitempty" json:"DepartureAirport,omitempty"`
	FareEquipment        *FareEquipment        `xml:" Equipment,omitempty" json:"Equipment,omitempty"`
	FareMarketingAirline *FareMarketingAirline `xml:" MarketingAirline,omitempty" json:"MarketingAirline,omitempty"`
}

type FareGroupType struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type FareMarketingAirline struct {
	AttrCode string `xml:" Code,attr"  json:",omitempty"`
}

type FareOTA_AirPriceRQ struct {
	// AttrEchoToken           string                   `xml:" EchoToken,attr"  json:",omitempty"`
	// AttrTarget              string                   `xml:" Target,attr"  json:",omitempty"`
	// AttrVersion             string                   `xml:" Version,attr"  json:",omitempty"`
	FareAirItinerary        *FareAirItinerary        `xml:" AirItinerary,omitempty" json:"AirItinerary,omitempty"`
	FareDistributor         *FareDistributor         `xml:" Distributor,omitempty" json:"Distributor,omitempty"`
	FareGroupType           *FareGroupType           `xml:" GroupType,omitempty" json:"GroupType,omitempty"`
	FarePOS                 *FarePOS                 `xml:" POS,omitempty" json:"POS,omitempty"`
	FareTravelerInfoSummary *FareTravelerInfoSummary `xml:" TravelerInfoSummary,omitempty" json:"TravelerInfoSummary,omitempty"`
}

type FareOriginDestinationOption struct {
	FareFlightSegment *FareFlightSegment `xml:" FlightSegment,omitempty" json:"FlightSegment,omitempty"`
}

type FareOriginDestinationOptions struct {
	FareOriginDestinationOption []*FareOriginDestinationOption `xml:" OriginDestinationOption,omitempty" json:"OriginDestinationOption,omitempty"`
}

type FarePOS struct {
	FareSource []*FareSource `xml:" Source,omitempty" json:"Source,omitempty"`
}

type FarePassengerTypeQuantity struct {
	AttrCode     string `xml:" Code,attr"  json:",omitempty"`
	AttrQuantity string `xml:" Quantity,attr"  json:",omitempty"`
}

type FarePriceRequestInformation struct {
	AttrCurrencyCode        string `xml:" CurrencyCode,attr"  json:",omitempty"`
	AttrNegotiatedFaresOnly string `xml:" NegotiatedFaresOnly,attr"  json:",omitempty"`
	AttrPricingSource       string `xml:" PricingSource,attr"  json:",omitempty"`
	AttrReprice             string `xml:" Reprice,attr"  json:",omitempty"`
}

type FareRequest struct {
	FareSITA_AirfarePriceRQ *FareSITA_AirfarePriceRQ `xml:" SITA_AirfarePriceRQ,omitempty" json:"SITA_AirfarePriceRQ,omitempty"`
}

type FareRequestorID struct {
	AttrID         string `xml:" ID,attr"  json:",omitempty"`
	AttrID_Context string `xml:" ID_Context,attr"  json:",omitempty"`
	AttrType       string `xml:" Type,attr"  json:",omitempty"`
}

type FareSITA_AirfarePriceRQ struct {
	//AttrVersion               string                     `xml:" Version,attr"  json:",omitempty"`
	FareAdditionalPriceRQData *FareAdditionalPriceRQData `xml:" AdditionalPriceRQData,omitempty" json:"AdditionalPriceRQData,omitempty"`
	FareOTA_AirPriceRQ        *FareOTA_AirPriceRQ        `xml:" OTA_AirPriceRQ,omitempty" json:"OTA_AirPriceRQ,omitempty"`
}

type FareSource struct {
	AttrAirportCode    string           `xml:" AirportCode,attr"  json:",omitempty"`
	AttrPseudoCityCode string           `xml:" PseudoCityCode,attr"  json:",omitempty"`
	FareRequestorID    *FareRequestorID `xml:" RequestorID,omitempty" json:"RequestorID,omitempty"`
}

type FareTSK_AirfarePrice struct {
	FareRequest *FareRequest `xml:" Request,omitempty" json:"Request,omitempty"`
	XMLName     xml.Name     `xml:"TSK_AirfarePrice"`
}

type FareTicketingCarrier struct {
	AttrCode string `xml:" Code,attr"  json:",omitempty"`
}

type FareTravelerInfoSummary struct {
	FareAirTravelerAvail        *FareAirTravelerAvail        `xml:" AirTravelerAvail,omitempty" json:"AirTravelerAvail,omitempty"`
	FarePriceRequestInformation *FarePriceRequestInformation `xml:" PriceRequestInformation,omitempty" json:"PriceRequestInformation,omitempty"`
}
