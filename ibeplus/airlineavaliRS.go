package ibeplus

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/otwdev/ibepluslib/models"

	"github.com/otwdev/galaxylib"

	"github.com/asaskevich/govalidator"

	golinq "github.com/ahmetb/go-linq"
)

type AirFarePriceResponse struct {
	PNR *models.PnrInfo
	RQ  *FareRoot
}

func NewFarePriceResponse(pnr *models.PnrInfo, rq *FareRoot) *AirFarePriceResponse {
	return &AirFarePriceResponse{
		PNR: pnr,
		RQ:  rq,
	}
}

const farePriceURL = "http://agibe.travelsky.com/ota/xml/AirFarePrice/I" // "http://ibeplus.travelsky.com/ota/xml/AirFarePricePolicy/I" //"http://agibe.travelsky.com/ota/xml/AirFarePrice/I"

func (a *AirFarePriceResponse) AirFarePriceRS() (*galaxylib.GalaxyError, int64) {

	ibe := NewIBE(farePriceURL, a.RQ.FareTSK_AirfarePrice)

	bytes, err := ibe.Reqeust() //(farePriceURL, a.RQ.FareTSK_AirfarePrice)

	if err != nil {
		return galaxylib.DefaultGalaxyError.FromError(1, err), 0
	}

	rev := string(bytes)
	fmt.Println(rev)

	airFarePrice := &RSFarePriceTSK_AirfarePrice{}

	if errInner := xml.Unmarshal(bytes, airFarePrice); errInner != nil {
		return galaxylib.DefaultGalaxyError.FromError(1, errInner), 0
	}

	//运价不匹配
	if airFarePrice.RSFarePriceResponse.RSFarePricePriceError != nil {
		return galaxylib.DefaultGalaxyError.FromText(3, airFarePrice.RSFarePriceResponse.RSFarePricePriceError.RSFarePriceCNMessage.Text), 0
	}
	if airFarePrice.RSFarePriceResponse.RSFarePriceSITA_AirfarePriceRS.RSFarePriceOTA_AirPriceRS.RSFarePriceWarnings != nil {
		return galaxylib.DefaultGalaxyError.FromText(3, "运价信息不匹配"), 0
	}
	if airFarePrice.RSFarePriceResponse.RSFarePriceSITA_AirfarePriceRS.RSFarePriceOTA_AirPriceRS.RSFarePricePricedItineraries.RSFarePricePolicyBindings == nil {
		galaxylib.GalaxyLogger.Debug("price:%s", a.RQ.FareTSK_AirfarePrice.FareRequest.FareSITA_AirfarePriceRQ.FareOTA_AirPriceRQ.FareAirItinerary.FareOriginDestinationOptions.FareOriginDestinationOption)
	} else {
		for _, v := range airFarePrice.RSFarePriceResponse.RSFarePriceSITA_AirfarePriceRS.RSFarePriceOTA_AirPriceRS.RSFarePricePricedItineraries.RSFarePricePolicyBindings.RSFarePricePolicys.RSFarePricePolicy {

			a.PNR.SegmentPrice = append(a.PNR.SegmentPrice, &models.FlightSegmentFarePrice{
				PNRID:    a.PNR.PnrID,
				JSONData: v.RSFarePriceContent.Text,
			})
		}
	}

	// //舱位
	if err, num := a.varifyClasses(airFarePrice); err != nil {
		return err, num
	}

	return nil, 0
}

func (a *AirFarePriceResponse) varifyClasses(rs *RSFarePriceTSK_AirfarePrice) (retErr *galaxylib.GalaxyError, combinNum int64) {
	// //舱位
	segmentsOpts := rs.RSFarePriceResponse.RSFarePriceSITA_AirfarePriceRS.RSFarePriceOTA_AirPriceRS.RSFarePricePricedItineraries.
		RSFarePricePricedItinerary[0].RSFarePriceAirItinerary.RSFarePriceOriginDestinationOptions.RSFarePriceOriginDestinationOption

	travelCount := a.PNR.ADTQuantity + a.PNR.CHDQuantity

	for _, v := range segmentsOpts {
		//舱位数量
		if err, num := verifyOrderCabin(v.RSFarePriceFlightSegment, a.PNR, travelCount); err != nil {
			return err, num //galaxylib.DefaultGalaxyError.FromError(1, err)
		}
	}

	aryPrice := rs.RSFarePriceResponse.RSFarePriceSITA_AirfarePriceRS.RSFarePriceOTA_AirPriceRS.RSFarePricePricedItineraries.RSFarePricePricedItinerary

	if err := a.getPrice(aryPrice); err != nil {
		return err, 0
	}

	// if a.PNR.ADTQuantity > 0 && a.PNR.SupplierADTPrice == 0 {
	// 	return utils.NewError(3, "价格异常")
	// }
	// if a.PNR.CHDQuantity > 0 && a.PNR.SupplierCHDPrice == 0 {
	// 	return utils.NewError(3, "价格异常")
	// }

	// if a.PNR.TotalSupplierPrice() > a.PNR.TotalCustomerPrice() {
	// 	logpkg.NewLogF("price", "价格匹配S-%d;C-%d", a.PNR.TotalSupplierPrice(), a.PNR.TotalCustomerPrice())
	// 	return utils.NewError(1, "不匹配的结算信")
	// }

	return nil, 0
}

func (a *AirFarePriceResponse) getPrice(aryPrice []*RSFarePricePricedItinerary) *galaxylib.GalaxyError {

	//revErr := &utils.NewError(3, "无价格信息")

	for _, v := range aryPrice {

		fareBasePrice := v.RSFarePriceAirItineraryPricingInfo.RSFarePricePTC_FareBreakdowns.RSFarePricePTC_FareBreakdown.RSFarePricePassengerFare.RSFarePriceBaseFare
		//票面价
		baseAmountStr := fareBasePrice.AttrAmount
		baseAmount, err := govalidator.ToFloat(baseAmountStr)
		if err != nil {
			return galaxylib.DefaultGalaxyError.FromText(1, fmt.Sprintf("票价格错误:%v", err))
		}

		//总价(含税)
		totalAmountStr := v.RSFarePriceAirItineraryPricingInfo.RSFarePricePTC_FareBreakdowns.RSFarePricePTC_FareBreakdown.RSFarePricePassengerFare.RSFarePriceTotalFare.AttrAmount
		totalAmount, err := govalidator.ToFloat(totalAmountStr)
		if err != nil {
			return galaxylib.DefaultGalaxyError.FromText(1, fmt.Sprintf("价格解决异常：%v", err))
		}
		//
		personType := v.RSFarePriceAirItineraryPricingInfo.RSFarePricePTC_FareBreakdowns.RSFarePricePTC_FareBreakdown.RSFarePricePassengerTypeQuantity.AttrCode

		rate, _ := govalidator.ToFloat(fareBasePrice.AttrRate)

		a.PNR.Price = append(a.PNR.Price, &models.PnrPrice{
			PNRID:        a.PNR.PnrID,
			DownPrice:    baseAmount,
			DownFax:      totalAmount - baseAmount,
			Type:         personType,
			ToCurrency:   fareBasePrice.AttrToCurrency,
			FromCurrency: fareBasePrice.AttrFromCurrency,
			Rate:         rate, //fareBasePrice.AttrRate,
		})

		// if personType == string(air.Adult) {
		// 	a.PNR.SupplierADTPrice = baseAmount
		// 	a.PNR.ADTFaxes = totalAmount - baseAmount
		// 	continue
		// }
		// if personType == string(air.Child) {
		// 	a.PNR.SupplierCHDPrice = baseAmount
		// 	a.PNR.CHDFaxes = totalAmount - baseAmount
		// }

	}

	return nil
}

// func verifycabin(travelCount int, opt *RSFarePriceOriginDestinationOption, cabin string) *utils.CustomError {
// 	aryClassAvail := opt.RSFarePriceFlightSegment.RSFarePriceBookingClassAvails.RSFarePriceBookingClassAvail
// 	for _, class := range aryClassAvail {
// 		if class.AttrResBookDesigCode == cabin {
// 			quantity, _ := strconv.ParseInt(class.AttrResBookDesigQuantity, 0, 10)
// 			if quantity < int64(travelCount) {
// 				return utils.NewError(5, "舱位不足")
// 			}
// 			return nil
// 		}
// 	}
// 	return utils.NewError(5, "无舱位")
// }

func verifyOrderCabin(segment *RSFarePriceFlightSegment, pnr *models.PnrInfo, travelCount int) (revErr *galaxylib.GalaxyError, combinNum int64) {

	//舱位名称
	cabin := golinq.From(pnr.FlightSegments).WhereT(func(f *models.FlightSegment) bool {
		return segment.AttrArrivalDateTime == f.ArriveDateTime() &&
			segment.AttrDepartureDateTime == f.DepartrueDateTime() &&
			segment.AttrFlightNumber == f.FlyNo &&
			segment.RSFarePriceMarketingAirline.AttrCode == f.MarketingAirLine &&
			segment.RSFarePriceArrivalAirport.AttrLocationCode == f.ArriveCityCode &&
			segment.RSFarePriceDepartureAirport.AttrLocationCode == f.DepartCityCode
	}).SelectT(func(f *models.FlightSegment) string {
		return f.Cabin
	}).First()
	strCabin := cabin.(string)

	if strCabin == "" {
		return galaxylib.DefaultGalaxyError.FromText(5, "无舱位"), 0
	}

	aryClassAvail := segment.RSFarePriceBookingClassAvails.RSFarePriceBookingClassAvail
	for _, class := range aryClassAvail {
		if class.AttrResBookDesigCode == cabin {
			quantity, _ := strconv.ParseInt(class.AttrResBookDesigQuantity, 0, 10)
			if quantity < int64(travelCount) {
				return galaxylib.DefaultGalaxyError.FromText(5, "舱位不足"), quantity
			}
			return nil, quantity
		}
	}
	return galaxylib.DefaultGalaxyError.FromText(5, "无舱位"), 0

}

// func getOrderCabin(segment *RSFarePriceFlightSegment, pnr *air.PnrInfo) string {
// 	cabin := golinq.From(pnr.FlightSegments).WhereT(func(f *air.FlightSegment) bool {
// 		return segment.AttrArrivalDateTime == f.ArriveDateTime() &&
// 			segment.AttrDepartureDateTime == f.DepartrueDateTime() &&
// 			segment.AttrFlightNumber == f.FlyNo &&
// 			segment.RSFarePriceMarketingAirline.AttrCode == f.MarketingAirLine
// 	}).SelectT(func(f *air.FlightSegment) string {
// 		return f.Cabin
// 	}).First()
// 	strCabin := cabin.(string)
// 	return strCabin
// }

type RSFarePriceRoot struct {
	RSFarePriceTSK_AirfarePrice *RSFarePriceTSK_AirfarePrice `xml:"http://www.travelsky.com/fare/xmlInterface TSK_AirfarePrice,omitempty" json:"TSK_AirfarePrice,omitempty"`
}

type RSFarePriceAirItinerary struct {
	RSFarePriceOriginDestinationOptions *RSFarePriceOriginDestinationOptions `xml:"http://www.opentravel.org/OTA/2003/05 OriginDestinationOptions,omitempty" json:"OriginDestinationOptions,omitempty"`
	XMLName                             xml.Name                             `xml:"http://www.opentravel.org/OTA/2003/05 AirItinerary,omitempty" json:"AirItinerary,omitempty"`
}

type RSFarePriceAirItineraryPricingInfo struct {
	RSFarePriceFareInfos          *RSFarePriceFareInfos          `xml:"http://www.opentravel.org/OTA/2003/05 FareInfos,omitempty" json:"FareInfos,omitempty"`
	RSFarePricePTC_FareBreakdowns *RSFarePricePTC_FareBreakdowns `xml:"http://www.opentravel.org/OTA/2003/05 PTC_FareBreakdowns,omitempty" json:"PTC_FareBreakdowns,omitempty"`
	XMLName                       xml.Name                       `xml:"http://www.opentravel.org/OTA/2003/05 AirItineraryPricingInfo,omitempty" json:"AirItineraryPricingInfo,omitempty"`
}

type RSFarePriceArrivalAirport struct {
	AttrCodeContext  string   `xml:" CodeContext,attr"  json:",omitempty"`
	AttrLocationCode string   `xml:" LocationCode,attr"  json:",omitempty"`
	XMLName          xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 ArrivalAirport,omitempty" json:"ArrivalAirport,omitempty"`
}

type RSFarePriceBaseFare struct {
	AttrAmount       string   `xml:" Amount,attr"  json:",omitempty"`
	AttrCurrencyCode string   `xml:" CurrencyCode,attr"  json:",omitempty"`
	AttrFromCurrency string   `xml:" FromCurrency,attr"  json:",omitempty"`
	AttrRate         string   `xml:" Rate,attr"  json:",omitempty"`
	AttrToCurrency   string   `xml:" ToCurrency,attr"  json:",omitempty"`
	XMLName          xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 BaseFare,omitempty" json:"BaseFare,omitempty"`
}

type RSFarePriceBookingClassAvail struct {
	AttrRPH                  string   `xml:" RPH,attr"  json:",omitempty"`
	AttrResBookDesigCode     string   `xml:" ResBookDesigCode,attr"  json:",omitempty"`
	AttrResBookDesigQuantity string   `xml:" ResBookDesigQuantity,attr"  json:",omitempty"`
	XMLName                  xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 BookingClassAvail,omitempty" json:"BookingClassAvail,omitempty"`
}

type RSFarePriceBookingClassAvails struct {
	RSFarePriceBookingClassAvail []*RSFarePriceBookingClassAvail `xml:"http://www.opentravel.org/OTA/2003/05 BookingClassAvail,omitempty" json:"BookingClassAvail,omitempty"`
	XMLName                      xml.Name                        `xml:"http://www.opentravel.org/OTA/2003/05 BookingClassAvails,omitempty" json:"BookingClassAvails,omitempty"`
}

type RSFarePriceDepartureAirport struct {
	AttrCodeContext  string   `xml:" CodeContext,attr"  json:",omitempty"`
	AttrLocationCode string   `xml:" LocationCode,attr"  json:",omitempty"`
	XMLName          xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 DepartureAirport,omitempty" json:"DepartureAirport,omitempty"`
}

type RSFarePriceFareBasisCode struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 FareBasisCode,omitempty" json:"FareBasisCode,omitempty"`
}

type RSFarePriceFareBasisCodes struct {
	RSFarePriceFareBasisCode []*RSFarePriceFareBasisCode `xml:"http://www.opentravel.org/OTA/2003/05 FareBasisCode,omitempty" json:"FareBasisCode,omitempty"`
	XMLName                  xml.Name                    `xml:"http://www.opentravel.org/OTA/2003/05 FareBasisCodes,omitempty" json:"FareBasisCodes,omitempty"`
}

type RSFarePriceFareInfo struct {
	AttrNegotiatedFare          string                       `xml:" NegotiatedFare,attr"  json:",omitempty"`
	RSFarePriceArrivalAirport   *RSFarePriceArrivalAirport   `xml:"http://www.opentravel.org/OTA/2003/05 ArrivalAirport,omitempty" json:"ArrivalAirport,omitempty"`
	RSFarePriceDepartureAirport *RSFarePriceDepartureAirport `xml:"http://www.opentravel.org/OTA/2003/05 DepartureAirport,omitempty" json:"DepartureAirport,omitempty"`
	RSFarePriceFareReference    *RSFarePriceFareReference    `xml:"http://www.opentravel.org/OTA/2003/05 FareReference,omitempty" json:"FareReference,omitempty"`
	RSFarePriceFilingAirline    *RSFarePriceFilingAirline    `xml:"http://www.opentravel.org/OTA/2003/05 FilingAirline,omitempty" json:"FilingAirline,omitempty"`
	RSFarePriceRuleInfo         *RSFarePriceRuleInfo         `xml:"http://www.opentravel.org/OTA/2003/05 RuleInfo,omitempty" json:"RuleInfo,omitempty"`
	RSFarePriceTPA_Extensions   *RSFarePriceTPA_Extensions   `xml:"http://www.opentravel.org/OTA/2003/05 TPA_Extensions,omitempty" json:"TPA_Extensions,omitempty"`
	XMLName                     xml.Name                     `xml:"http://www.opentravel.org/OTA/2003/05 FareInfo,omitempty" json:"FareInfo,omitempty"`
}

type RSFarePriceFareInfos struct {
	RSFarePriceFareInfo []*RSFarePriceFareInfo `xml:"http://www.opentravel.org/OTA/2003/05 FareInfo,omitempty" json:"FareInfo,omitempty"`
	XMLName             xml.Name               `xml:"http://www.opentravel.org/OTA/2003/05 FareInfos,omitempty" json:"FareInfos,omitempty"`
}

type RSFarePriceFareReference struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 FareReference,omitempty" json:"FareReference,omitempty"`
}

type RSFarePriceFilingAirline struct {
	AttrCode string   `xml:" Code,attr"  json:",omitempty"`
	XMLName  xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 FilingAirline,omitempty" json:"FilingAirline,omitempty"`
}

type RSFarePriceFlightSegment struct {
	AttrArrivalDateTime           string                         `xml:" ArrivalDateTime,attr"  json:",omitempty"`
	AttrDepartureDateTime         string                         `xml:" DepartureDateTime,attr"  json:",omitempty"`
	AttrFlightNumber              string                         `xml:" FlightNumber,attr"  json:",omitempty"`
	AttrRPH                       string                         `xml:" RPH,attr"  json:",omitempty"`
	AttrResBookDesigCode          string                         `xml:" ResBookDesigCode,attr"  json:",omitempty"`
	RSFarePriceArrivalAirport     *RSFarePriceArrivalAirport     `xml:"http://www.opentravel.org/OTA/2003/05 ArrivalAirport,omitempty" json:"ArrivalAirport,omitempty"`
	RSFarePriceBookingClassAvails *RSFarePriceBookingClassAvails `xml:"http://www.opentravel.org/OTA/2003/05 BookingClassAvails,omitempty" json:"BookingClassAvails,omitempty"`
	RSFarePriceDepartureAirport   *RSFarePriceDepartureAirport   `xml:"http://www.opentravel.org/OTA/2003/05 DepartureAirport,omitempty" json:"DepartureAirport,omitempty"`
	RSFarePriceMarketingAirline   *RSFarePriceMarketingAirline   `xml:"http://www.opentravel.org/OTA/2003/05 MarketingAirline,omitempty" json:"MarketingAirline,omitempty"`
	XMLName                       xml.Name                       `xml:"http://www.opentravel.org/OTA/2003/05 FlightSegment,omitempty" json:"FlightSegment,omitempty"`
}

type RSFarePriceMarketingAirline struct {
	AttrCode string   `xml:" Code,attr"  json:",omitempty"`
	XMLName  xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 MarketingAirline,omitempty" json:"MarketingAirline,omitempty"`
}

type RSFarePriceNotes struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 Notes,omitempty" json:"Notes,omitempty"`
}

type RSFarePriceOTA_AirPriceRS struct {
	AttrTarget                   string                        `xml:" Target,attr"  json:",omitempty"`
	AttrTimeStamp                string                        `xml:" TimeStamp,attr"  json:",omitempty"`
	AttrVersion                  string                        `xml:" Version,attr"  json:",omitempty"`
	RSFarePricePricedItineraries *RSFarePricePricedItineraries `xml:"http://www.opentravel.org/OTA/2003/05 PricedItineraries,omitempty" json:"PricedItineraries,omitempty"`
	RSFarePriceSuccess           *RSFarePriceSuccess           `xml:"http://www.opentravel.org/OTA/2003/05 Success,omitempty" json:"Success,omitempty"`
	RSFarePriceWarnings          *RSFarePriceWarnings          `xml:"http://www.opentravel.org/OTA/2003/05 Warnings,omitempty" json:"Warnings,omitempty"`
	XMLName                      xml.Name                      `xml:"http://www.opentravel.org/OTA/2003/05 OTA_AirPriceRS,omitempty" json:"OTA_AirPriceRS,omitempty"`
}

type RSFarePriceOriginDestinationOption struct {
	RSFarePriceFlightSegment *RSFarePriceFlightSegment `xml:"http://www.opentravel.org/OTA/2003/05 FlightSegment,omitempty" json:"FlightSegment,omitempty"`
	XMLName                  xml.Name                  `xml:"http://www.opentravel.org/OTA/2003/05 OriginDestinationOption,omitempty" json:"OriginDestinationOption,omitempty"`
}

type RSFarePriceOriginDestinationOptions struct {
	RSFarePriceOriginDestinationOption []*RSFarePriceOriginDestinationOption `xml:"http://www.opentravel.org/OTA/2003/05 OriginDestinationOption,omitempty" json:"OriginDestinationOption,omitempty"`
	XMLName                            xml.Name                              `xml:"http://www.opentravel.org/OTA/2003/05 OriginDestinationOptions,omitempty" json:"OriginDestinationOptions,omitempty"`
}

type RSFarePricePTC_FareBreakdown struct {
	AttrPricingSource                string                            `xml:" PricingSource,attr"  json:",omitempty"`
	RSFarePriceFareBasisCodes        *RSFarePriceFareBasisCodes        `xml:"http://www.opentravel.org/OTA/2003/05 FareBasisCodes,omitempty" json:"FareBasisCodes,omitempty"`
	RSFarePricePassengerFare         *RSFarePricePassengerFare         `xml:"http://www.opentravel.org/OTA/2003/05 PassengerFare,omitempty" json:"PassengerFare,omitempty"`
	RSFarePricePassengerTypeQuantity *RSFarePricePassengerTypeQuantity `xml:"http://www.opentravel.org/OTA/2003/05 PassengerTypeQuantity,omitempty" json:"PassengerTypeQuantity,omitempty"`
	XMLName                          xml.Name                          `xml:"http://www.opentravel.org/OTA/2003/05 PTC_FareBreakdown,omitempty" json:"PTC_FareBreakdown,omitempty"`
}

type RSFarePricePTC_FareBreakdowns struct {
	RSFarePricePTC_FareBreakdown *RSFarePricePTC_FareBreakdown `xml:"http://www.opentravel.org/OTA/2003/05 PTC_FareBreakdown,omitempty" json:"PTC_FareBreakdown,omitempty"`
	XMLName                      xml.Name                      `xml:"http://www.opentravel.org/OTA/2003/05 PTC_FareBreakdowns,omitempty" json:"PTC_FareBreakdowns,omitempty"`
}

type RSFarePricePassengerFare struct {
	RSFarePriceBaseFare             *RSFarePriceBaseFare             `xml:"http://www.opentravel.org/OTA/2003/05 BaseFare,omitempty" json:"BaseFare,omitempty"`
	RSFarePriceTPA_Extensions       *RSFarePriceTPA_Extensions       `xml:"http://www.opentravel.org/OTA/2003/05 TPA_Extensions,omitempty" json:"TPA_Extensions,omitempty"`
	RSFarePriceTaxes                *RSFarePriceTaxes                `xml:"http://www.opentravel.org/OTA/2003/05 Taxes,omitempty" json:"Taxes,omitempty"`
	RSFarePriceTotalFare            *RSFarePriceTotalFare            `xml:"http://www.opentravel.org/OTA/2003/05 TotalFare,omitempty" json:"TotalFare,omitempty"`
	RSFarePriceUnstructuredFareCalc *RSFarePriceUnstructuredFareCalc `xml:"http://www.opentravel.org/OTA/2003/05 UnstructuredFareCalc,omitempty" json:"UnstructuredFareCalc,omitempty"`
	XMLName                         xml.Name                         `xml:"http://www.opentravel.org/OTA/2003/05 PassengerFare,omitempty" json:"PassengerFare,omitempty"`
}

type RSFarePricePassengerTypeQuantity struct {
	AttrCode string   `xml:" Code,attr"  json:",omitempty"`
	XMLName  xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 PassengerTypeQuantity,omitempty" json:"PassengerTypeQuantity,omitempty"`
}

type RSFarePricePricedItineraries struct {
	RSFarePricePolicyBindings  *RSFarePricePolicyBindings    `xml:"http://www.travelsky.com/fare/xmlInterface PolicyBindings,omitempty" json:"PolicyBindings,omitempty"`
	RSFarePricePricedItinerary []*RSFarePricePricedItinerary `xml:"http://www.opentravel.org/OTA/2003/05 PricedItinerary,omitempty" json:"PricedItinerary,omitempty"`
	XMLName                    xml.Name                      `xml:"http://www.opentravel.org/OTA/2003/05 PricedItineraries,omitempty" json:"PricedItineraries,omitempty"`
}

type RSFarePricePricedItinerary struct {
	AttrSequenceNumber                 string                              `xml:" SequenceNumber,attr"  json:",omitempty"`
	RSFarePriceAirItinerary            *RSFarePriceAirItinerary            `xml:"http://www.opentravel.org/OTA/2003/05 AirItinerary,omitempty" json:"AirItinerary,omitempty"`
	RSFarePriceAirItineraryPricingInfo *RSFarePriceAirItineraryPricingInfo `xml:"http://www.opentravel.org/OTA/2003/05 AirItineraryPricingInfo,omitempty" json:"AirItineraryPricingInfo,omitempty"`
	RSFarePriceNotes                   []*RSFarePriceNotes                 `xml:"http://www.opentravel.org/OTA/2003/05 Notes,omitempty" json:"Notes,omitempty"`
	RSFarePriceTicketingInfo           *RSFarePriceTicketingInfo           `xml:"http://www.opentravel.org/OTA/2003/05 TicketingInfo,omitempty" json:"TicketingInfo,omitempty"`
	XMLName                            xml.Name                            `xml:"http://www.opentravel.org/OTA/2003/05 PricedItinerary,omitempty" json:"PricedItinerary,omitempty"`
}

type RSFarePriceRuleInfo struct {
	XMLName xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 RuleInfo,omitempty" json:"RuleInfo,omitempty"`
}

type RSFarePriceSITA_PassengerFareExtension struct {
	RSFarePriceAppliedPTCs *RSFarePriceAppliedPTCs `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS AppliedPTCs,omitempty" json:"AppliedPTCs,omitempty"`
	XMLName                xml.Name                `xml:"http://www.opentravel.org/OTA/2003/05 SITA_PassengerFareExtension,omitempty" json:"SITA_PassengerFareExtension,omitempty"`
}

type RSFarePriceSuccess struct {
	XMLName xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 Success,omitempty" json:"Success,omitempty"`
}

type RSFarePriceTPA_Extensions struct {
	RSFarePriceSITA_FareInfoExtension      *RSFarePriceSITA_FareInfoExtension      `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS SITA_FareInfoExtension,omitempty" json:"SITA_FareInfoExtension,omitempty"`
	RSFarePriceSITA_PassengerFareExtension *RSFarePriceSITA_PassengerFareExtension `xml:"http://www.opentravel.org/OTA/2003/05 SITA_PassengerFareExtension,omitempty" json:"SITA_PassengerFareExtension,omitempty"`
	XMLName                                xml.Name                                `xml:"http://www.opentravel.org/OTA/2003/05 TPA_Extensions,omitempty" json:"TPA_Extensions,omitempty"`
}

type RSFarePriceTax struct {
	AttrAmount       string   `xml:" Amount,attr"  json:",omitempty"`
	AttrCurrencyCode string   `xml:" CurrencyCode,attr"  json:",omitempty"`
	AttrTaxCode      string   `xml:" TaxCode,attr"  json:",omitempty"`
	XMLName          xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 Tax,omitempty" json:"Tax,omitempty"`
}

type RSFarePriceTaxes struct {
	RSFarePriceTax []*RSFarePriceTax `xml:"http://www.opentravel.org/OTA/2003/05 Tax,omitempty" json:"Tax,omitempty"`
	XMLName        xml.Name          `xml:"http://www.opentravel.org/OTA/2003/05 Taxes,omitempty" json:"Taxes,omitempty"`
}

type RSFarePriceTicketingInfo struct {
	AttrTicketTimeLimit string   `xml:" TicketTimeLimit,attr"  json:",omitempty"`
	AttrTicketType      string   `xml:" TicketType,attr"  json:",omitempty"`
	XMLName             xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 TicketingInfo,omitempty" json:"TicketingInfo,omitempty"`
}

type RSFarePriceTotalFare struct {
	AttrAmount       string   `xml:" Amount,attr"  json:",omitempty"`
	AttrCurrencyCode string   `xml:" CurrencyCode,attr"  json:",omitempty"`
	XMLName          xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 TotalFare,omitempty" json:"TotalFare,omitempty"`
}

type RSFarePriceUnstructuredFareCalc struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 UnstructuredFareCalc,omitempty" json:"UnstructuredFareCalc,omitempty"`
}

type RSFarePriceWarning struct {
	AttrType string   `xml:" Type,attr"  json:",omitempty"`
	Text     string   `xml:",chardata" json:",omitempty"`
	XMLName  xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 Warning,omitempty" json:"Warning,omitempty"`
}

type RSFarePriceWarnings struct {
	RSFarePriceWarning []*RSFarePriceWarning `xml:"http://www.opentravel.org/OTA/2003/05 Warning,omitempty" json:"Warning,omitempty"`
	XMLName            xml.Name              `xml:"http://www.opentravel.org/OTA/2003/05 Warnings,omitempty" json:"Warnings,omitempty"`
}

type RSFarePriceAdditionalItinerariesData struct {
	RSFarePriceAdditionalItineraryData []*RSFarePriceAdditionalItineraryData `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS AdditionalItineraryData,omitempty" json:"AdditionalItineraryData,omitempty"`
	XMLName                            xml.Name                              `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS AdditionalItinerariesData,omitempty" json:"AdditionalItinerariesData,omitempty"`
}

type RSFarePriceAdditionalItineraryData struct {
	AttrFlightRestrictionInd          string                             `xml:" FlightRestrictionInd,attr"  json:",omitempty"`
	AttrFunctionalControlCharacter    string                             `xml:" FunctionalControlCharacter,attr"  json:",omitempty"`
	AttrPublicPrivateConstructInd     string                             `xml:" PublicPrivateConstructInd,attr"  json:",omitempty"`
	AttrRBD_OverrideInd               string                             `xml:" RBD_OverrideInd,attr"  json:",omitempty"`
	AttrRBD_RestrictionInd            string                             `xml:" RBD_RestrictionInd,attr"  json:",omitempty"`
	AttrSequenceNumber                string                             `xml:" SequenceNumber,attr"  json:",omitempty"`
	RSFarePriceAdditionalSegmentInfos *RSFarePriceAdditionalSegmentInfos `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS AdditionalSegmentInfos,omitempty" json:"AdditionalSegmentInfos,omitempty"`
	RSFarePriceAdditionalTicketInfo   *RSFarePriceAdditionalTicketInfo   `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS AdditionalTicketInfo,omitempty" json:"AdditionalTicketInfo,omitempty"`
	XMLName                           xml.Name                           `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS AdditionalItineraryData,omitempty" json:"AdditionalItineraryData,omitempty"`
}

type RSFarePriceAdditionalPriceRSData struct {
	RSFarePriceAdditionalItinerariesData *RSFarePriceAdditionalItinerariesData `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS AdditionalItinerariesData,omitempty" json:"AdditionalItinerariesData,omitempty"`
	XMLName                              xml.Name                              `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS AdditionalPriceRSData,omitempty" json:"AdditionalPriceRSData,omitempty"`
}

type RSFarePriceAdditionalSegmentInfo struct {
	AttrFareBasis                      string                              `xml:" FareBasis,attr"  json:",omitempty"`
	AttrFareRPH                        string                              `xml:" FareRPH,attr"  json:",omitempty"`
	AttrFreeBaggageAllowance           string                              `xml:" FreeBaggageAllowance,attr"  json:",omitempty"`
	AttrGlobalIndicator                string                              `xml:" GlobalIndicator,attr"  json:",omitempty"`
	AttrNotValidAfter                  string                              `xml:" NotValidAfter,attr"  json:",omitempty"`
	AttrNotValidBefore                 string                              `xml:" NotValidBefore,attr"  json:",omitempty"`
	AttrSegmentRPH                     string                              `xml:" SegmentRPH,attr"  json:",omitempty"`
	AttrStopoverPermitted              string                              `xml:" StopoverPermitted,attr"  json:",omitempty"`
	RSFarePriceRebookResBookDesigCodes *RSFarePriceRebookResBookDesigCodes `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS RebookResBookDesigCodes,omitempty" json:"RebookResBookDesigCodes,omitempty"`
	XMLName                            xml.Name                            `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS AdditionalSegmentInfo,omitempty" json:"AdditionalSegmentInfo,omitempty"`
}

type RSFarePriceAdditionalSegmentInfos struct {
	RSFarePriceAdditionalSegmentInfo []*RSFarePriceAdditionalSegmentInfo `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS AdditionalSegmentInfo,omitempty" json:"AdditionalSegmentInfo,omitempty"`
	XMLName                          xml.Name                            `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS AdditionalSegmentInfos,omitempty" json:"AdditionalSegmentInfos,omitempty"`
}

type RSFarePriceAdditionalTicketInfo struct {
	RSFarePriceEndorsements *RSFarePriceEndorsements `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS Endorsements,omitempty" json:"Endorsements,omitempty"`
	XMLName                 xml.Name                 `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS AdditionalTicketInfo,omitempty" json:"AdditionalTicketInfo,omitempty"`
}

type RSFarePriceAppliedPTCs struct {
	AttrPTC string   `xml:" PTC,attr"  json:",omitempty"`
	XMLName xml.Name `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS AppliedPTCs,omitempty" json:"AppliedPTCs,omitempty"`
}

type RSFarePriceDirectionality struct {
	AttrCode string   `xml:" Code,attr"  json:",omitempty"`
	XMLName  xml.Name `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS Directionality,omitempty" json:"Directionality,omitempty"`
}

type RSFarePriceEndorsement struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS Endorsement,omitempty" json:"Endorsement,omitempty"`
}

type RSFarePriceEndorsements struct {
	RSFarePriceEndorsement *RSFarePriceEndorsement `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS Endorsement,omitempty" json:"Endorsement,omitempty"`
	XMLName                xml.Name                `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS Endorsements,omitempty" json:"Endorsements,omitempty"`
}

type RSFarePriceRebookResBookDesigCode struct {
	AttrResBookDesigCode string   `xml:" ResBookDesigCode,attr"  json:",omitempty"`
	XMLName              xml.Name `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS RebookResBookDesigCode,omitempty" json:"RebookResBookDesigCode,omitempty"`
}

type RSFarePriceRebookResBookDesigCodes struct {
	RSFarePriceRebookResBookDesigCode *RSFarePriceRebookResBookDesigCode `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS RebookResBookDesigCode,omitempty" json:"RebookResBookDesigCode,omitempty"`
	XMLName                           xml.Name                           `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS RebookResBookDesigCodes,omitempty" json:"RebookResBookDesigCodes,omitempty"`
}

type RSFarePriceRef1 struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS Ref1,omitempty" json:"Ref1,omitempty"`
}

type RSFarePriceRef2 struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS Ref2,omitempty" json:"Ref2,omitempty"`
}

type RSFarePriceReferences struct {
	RSFarePriceRef1 *RSFarePriceRef1 `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS Ref1,omitempty" json:"Ref1,omitempty"`
	RSFarePriceRef2 *RSFarePriceRef2 `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS Ref2,omitempty" json:"Ref2,omitempty"`
	XMLName         xml.Name         `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS References,omitempty" json:"References,omitempty"`
}

type RSFarePriceSITA_AirfarePriceRS struct {
	AttrXmlnsOta                     string                            `xml:"xmlns ota,attr"  json:",omitempty"`
	AttrXmlnsSita                    string                            `xml:"xmlns sita,attr"  json:",omitempty"`
	AttrVersion                      string                            `xml:" Version,attr"  json:",omitempty"`
	AttrXmlnsXsi                     string                            `xml:"xmlns xsi,attr"  json:",omitempty"`
	RSFarePriceAdditionalPriceRSData *RSFarePriceAdditionalPriceRSData `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS AdditionalPriceRSData,omitempty" json:"AdditionalPriceRSData,omitempty"`
	RSFarePriceOTA_AirPriceRS        *RSFarePriceOTA_AirPriceRS        `xml:"http://www.opentravel.org/OTA/2003/05 OTA_AirPriceRS,omitempty" json:"OTA_AirPriceRS,omitempty"`
	XMLName                          xml.Name                          `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS SITA_AirfarePriceRS,omitempty" json:"SITA_AirfarePriceRS,omitempty"`
}

type RSFarePriceSITA_FareInfoExtension struct {
	AttrFareComponentDirection       string                            `xml:" FareComponentDirection,attr"  json:",omitempty"`
	AttrFareRPH                      string                            `xml:" FareRPH,attr"  json:",omitempty"`
	AttrFareType                     string                            `xml:" FareType,attr"  json:",omitempty"`
	AttrFootNoteDate                 string                            `xml:" FootNoteDate,attr"  json:",omitempty"`
	AttrRuleNumber                   string                            `xml:" RuleNumber,attr"  json:",omitempty"`
	AttrRuleTariffNumber             string                            `xml:" RuleTariffNumber,attr"  json:",omitempty"`
	AttrTariffNumber                 string                            `xml:" TariffNumber,attr"  json:",omitempty"`
	RSFarePriceDirectionality        *RSFarePriceDirectionality        `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS Directionality,omitempty" json:"Directionality,omitempty"`
	RSFarePriceReferences            *RSFarePriceReferences            `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS References,omitempty" json:"References,omitempty"`
	RSFarePriceSubjectToGovtApproval *RSFarePriceSubjectToGovtApproval `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS SubjectToGovtApproval,omitempty" json:"SubjectToGovtApproval,omitempty"`
	XMLName                          xml.Name                          `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS SITA_FareInfoExtension,omitempty" json:"SITA_FareInfoExtension,omitempty"`
}

type RSFarePriceSubjectToGovtApproval struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS SubjectToGovtApproval,omitempty" json:"SubjectToGovtApproval,omitempty"`
}

type RSFarePriceCNMessage struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://www.travelsky.com/fare/xmlInterface CNMessage,omitempty" json:"CNMessage,omitempty"`
}

type RSFarePriceENMessage struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://www.travelsky.com/fare/xmlInterface ENMessage,omitempty" json:"ENMessage,omitempty"`
}

type RSFarePriceErrorNo struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://www.travelsky.com/fare/xmlInterface ErrorNo,omitempty" json:"ErrorNo,omitempty"`
}

type RSFarePricePolicy struct {
	RSFarePriceContent *RSFarePriceContent `xml:"http://www.travelsky.com/fare/xmlInterface content,omitempty" json:"content,omitempty"`
	RSFarePriceId      *RSFarePriceId      `xml:"http://www.travelsky.com/fare/xmlInterface id,omitempty" json:"id,omitempty"`
	RSFarePriceSeq     *RSFarePriceSeq     `xml:"http://www.travelsky.com/fare/xmlInterface seq,omitempty" json:"seq,omitempty"`
	XMLName            xml.Name            `xml:"http://www.travelsky.com/fare/xmlInterface Policy,omitempty" json:"Policy,omitempty"`
}

type RSFarePricePolicyBindings struct {
	RSFarePricePolicys *RSFarePricePolicys `xml:"http://www.travelsky.com/fare/xmlInterface Policys,omitempty" json:"Policys,omitempty"`
	XMLName            xml.Name            `xml:"http://www.travelsky.com/fare/xmlInterface PolicyBindings,omitempty" json:"PolicyBindings,omitempty"`
}

type RSFarePricePolicys struct {
	RSFarePricePolicy []*RSFarePricePolicy `xml:"http://www.travelsky.com/fare/xmlInterface Policy,omitempty" json:"Policy,omitempty"`
	XMLName           xml.Name             `xml:"http://www.travelsky.com/fare/xmlInterface Policys,omitempty" json:"Policys,omitempty"`
}

type RSFarePricePriceError struct {
	RSFarePriceCNMessage *RSFarePriceCNMessage `xml:"http://www.travelsky.com/fare/xmlInterface CNMessage,omitempty" json:"CNMessage,omitempty"`
	RSFarePriceENMessage *RSFarePriceENMessage `xml:"http://www.travelsky.com/fare/xmlInterface ENMessage,omitempty" json:"ENMessage,omitempty"`
	RSFarePriceErrorNo   *RSFarePriceErrorNo   `xml:"http://www.travelsky.com/fare/xmlInterface ErrorNo,omitempty" json:"ErrorNo,omitempty"`
	XMLName              xml.Name              `xml:"http://www.travelsky.com/fare/xmlInterface PriceError,omitempty" json:"PriceError,omitempty"`
}

type RSFarePriceResponse struct {
	RSFarePricePriceError          *RSFarePricePriceError          `xml:"http://www.travelsky.com/fare/xmlInterface PriceError,omitempty" json:"PriceError,omitempty"`
	RSFarePriceSITA_AirfarePriceRS *RSFarePriceSITA_AirfarePriceRS `xml:"http://www.sita.aero/PTS/fare/2005/11/PriceRS SITA_AirfarePriceRS,omitempty" json:"SITA_AirfarePriceRS,omitempty"`
	RSFarePriceToken               *RSFarePriceToken               `xml:"http://www.travelsky.com/fare/xmlInterface Token,omitempty" json:"Token,omitempty"`
	RSFarePriceVersion             *RSFarePriceVersion             `xml:"http://www.travelsky.com/fare/xmlInterface Version,omitempty" json:"Version,omitempty"`
	XMLName                        xml.Name                        `xml:"http://www.travelsky.com/fare/xmlInterface Response,omitempty" json:"Response,omitempty"`
}

type RSFarePriceTSK_AirfarePrice struct {
	AttrXmlns           string               `xml:" xmlns,attr"  json:",omitempty"`
	RSFarePriceResponse *RSFarePriceResponse `xml:"http://www.travelsky.com/fare/xmlInterface Response,omitempty" json:"Response,omitempty"`
	XMLName             xml.Name             `xml:"http://www.travelsky.com/fare/xmlInterface TSK_AirfarePrice,omitempty" json:"TSK_AirfarePrice,omitempty"`
}

type RSFarePriceToken struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://www.travelsky.com/fare/xmlInterface Token,omitempty" json:"Token,omitempty"`
}

type RSFarePriceVersion struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://www.travelsky.com/fare/xmlInterface Version,omitempty" json:"Version,omitempty"`
}

type RSFarePriceContent struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://www.travelsky.com/fare/xmlInterface content,omitempty" json:"content,omitempty"`
}

type RSFarePriceId struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://www.travelsky.com/fare/xmlInterface id,omitempty" json:"id,omitempty"`
}

type RSFarePriceSeq struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"http://www.travelsky.com/fare/xmlInterface seq,omitempty" json:"seq,omitempty"`
}
