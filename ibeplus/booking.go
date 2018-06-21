package ibeplus

import (
	"encoding/xml"
	"fmt"
	"ibepluslib/models"
	"time"

	"github.com/otwdev/galaxylib"

	"github.com/asaskevich/govalidator"
)

const bookingURL = "http://ibeplus.travelsky.com/ota/xml/AirBook"

type PNRBooking struct {
	Order *models.OrderInfo
	PNR   string
}

func NewPNRBooking(order *models.OrderInfo) *PNRBooking {
	return &PNRBooking{order, ""}
}

func (p *PNRBooking) makeBookingPerson(trvaler *models.Traveler) *BookingPersonName {

	if trvaler.IDCardType == string(models.IDCard) {
		return &BookingPersonName{
			AttrLanguageType: "ZH",
			BookingSurname: &BookingSurname{
				Text: trvaler.PersonName,
			},
		}
	}
	return &BookingPersonName{
		BookingSurname: &BookingSurname{
			Text: fmt.Sprintf("%s/%s", trvaler.GivenName, trvaler.SurName),
		},
		AttrLanguageType: "EN",
	}
}

func (p *PNRBooking) Booking() *galaxylib.GalaxyError {

	bk := &BookingRoot{}

	pnr := p.Order.PnrInofs[0]

	bk.BookingOTA_AirBookRQ = &BookingOTA_AirBookRQ{}

	//Office号
	bk.BookingOTA_AirBookRQ.BookingPOS = &BookingPOS{
		BookingSource: &BookingSource{
			AttrPseudoCityCode: pnr.OfficeNumber,
		},
	}

	//航段
	bk.BookingOTA_AirBookRQ.BookingAirItinerary = &BookingAirItinerary{}
	bk.BookingOTA_AirBookRQ.BookingAirItinerary.BookingOriginDestinationOptions = &BookingOriginDestinationOptions{}
	bk.BookingOTA_AirBookRQ.BookingAirItinerary.BookingOriginDestinationOptions.BookingOriginDestinationOption = &BookingOriginDestinationOption{}

	var segments []*BookingFlightSegment
	defaultAirline := ""

	for _, v := range pnr.FlightSegments {
		segments = append(segments, &BookingFlightSegment{
			AttrArrivalDateTime:   v.ArriveDateTime(),
			AttrDepartureDateTime: v.DepartrueDateTime(),
			AttrCodeshareInd:      "false",
			AttrFlighttNumber:     v.FlyNo,
			AttrStatus:            "NN",
			AttrSegmentType:       "NORMAL",
			AttrRPH:               govalidator.ToString(v.TripSeq),
			BookingDepartureAirport: &BookingDepartureAirport{
				AttrLocationCode: v.DepartCityCode,
			},
			BookingArrivalAirport: &BookingArrivalAirport{
				AttrLocationCode: v.ArriveCityCode,
			},
			BookingMarketingAirline: &BookingMarketingAirline{
				AttrCode: v.MarketingAirLine,
			},
			BookingBookingClassAvail: &BookingBookingClassAvail{
				AttrResBookDesigCode: v.Cabin,
			},
		})
		if defaultAirline == "" {
			defaultAirline = v.MarketingAirLine
		}
	}
	bk.BookingOTA_AirBookRQ.BookingAirItinerary.BookingOriginDestinationOptions.BookingOriginDestinationOption.BookingFlightSegment = segments

	//乘客
	bk.BookingOTA_AirBookRQ.BookingTravelerInfo = &BookingTravelerInfo{}
	var bookingTraveler []*BookingAirTraveler

	ctcm := ""

	for i, t := range pnr.TravelerInfos {
		trl := &BookingAirTraveler{}
		trl.AttrPassengerTypeCode = t.Type
		trl.AttrGender = t.Gender
		trl.BookingPersonName = p.makeBookingPerson(t)
		doc := &BookingDocument{}
		// doc.AttrDocHolderInd = "true"
		doc.AttrBirthDate = t.Birthday
		doc.AttrDocType = t.IDCardType
		doc.AttrDocID = t.IDCardNo
		doc.AttrDocHolderNationality = t.Nationality
		doc.AttrDocIssueCountry = t.IDIssueCountry
		doc.AttrBirthDate = t.Birthday
		doc.AttrGender = t.Gender
		doc.AttrExpireDate = t.IDExpireDate
		doc.AttrDocTypeDetail = "P"
		doc.AttrRPH = govalidator.ToString(i + 1)
		doc.BookingDocHolderFormattedName = &BookingDocHolderFormattedName{
			BookingGivenName: &BookingGivenName{
				Text: t.GivenName,
			},
			BookingSurname: &BookingSurname{
				Text: t.SurName,
			},
		}
		trl.BookingDocument = doc
		//rph,_ := govalidator.ToString()
		trl.BookingTravelerRefNumber = append(trl.BookingTravelerRefNumber, &BookingTravelerRefNumber{
			AttrRPH: govalidator.ToString(i + 1),
		})
		trl.BookingComment = &BookingComment{
			Text: "HK",
		}
		trl.BookingFlightSegmentRPHs = &BookingFlightSegmentRPHs{
			&BookingFlightSegmentRPH{"1"},
		}
		// trl.BookingDocumentFlightBinding = &BookingDocumentFlightBinding{
		// 	BookingDocumentRPH:      &BookingDocumentRPH{"1"},
		// 	BookingFlightSegmentRPH: &BookingFlightSegmentRPH{"1"},
		// }

		bookingTraveler = append(bookingTraveler, trl)

		if ctcm == "" {
			ctcm = t.Mobile
		}
	}

	bk.BookingOTA_AirBookRQ.BookingTravelerInfo.BookingAirTraveler = bookingTraveler
	bk.BookingOTA_AirBookRQ.BookingTravelerInfo.BookingSpecialReqDetails = &BookingSpecialReqDetails{}
	bk.BookingOTA_AirBookRQ.BookingTravelerInfo.BookingSpecialReqDetails.BookingOtherServiceInformations = &BookingOtherServiceInformations{}
	var oths []*BookingOtherServiceInformation
	oths = append(oths, &BookingOtherServiceInformation{
		AttrCode: "OTHS",
		BookingText: &BookingText{
			Text: fmt.Sprintf("CTCT%s", p.Order.ContactInfo.MobilePhone),
		},
		BookingAirline: &BookingAirline{
			AttrCode: defaultAirline,
		},
	})
	oths = append(oths, &BookingOtherServiceInformation{
		AttrCode: "OTHS",
		BookingText: &BookingText{
			Text: fmt.Sprintf("CTCM%s", ctcm),
		},
		BookingAirline: &BookingAirline{
			AttrCode: defaultAirline,
		},
		BookingTravelerRefNumber: []*BookingTravelerRefNumber{
			&BookingTravelerRefNumber{"1"},
		},
	})

	bk.BookingOTA_AirBookRQ.BookingTravelerInfo.BookingSpecialReqDetails.BookingOtherServiceInformations.BookingOtherServiceInformation = oths

	//出票信息

	depTime := pnr.FlightSegments[0].DepartrueDateTime()
	layout := "2006-01-02T15:04:05"
	limitTime, _ := time.Parse(layout, depTime)
	limitTime = limitTime.Add(-2 * time.Hour)
	//contract := p.Order.ContactInfo[0]

	bk.BookingOTA_AirBookRQ.BookingTicketing = &BookingTicketing{}
	bk.BookingOTA_AirBookRQ.BookingTicketing.AttrTicketTimeLimit = limitTime.Format("2006-01-02T15:04:05")
	bk.BookingOTA_AirBookRQ.BookingTPA_Extensions = &BookingTPA_Extensions{
		BookingContactInfo: &BookingContactInfo{
			Text: p.Order.ContactInfo.MobilePhone, //"13910556253",
		},
		BookingEnvelopType: &BookingEnvelopType{
			Text: "KI",
		},
	}

	if booking := galaxylib.GalaxyCfgFile.MustBool("booking", "enableBooking"); booking == false {
		p.PNR = "TEST1231"
		return nil
	}

	ibe := NewIBE(bookingURL, bk.BookingOTA_AirBookRQ)
	rev, err := ibe.Reqeust() //ReqeustIBE(bookingURL, bk.BookingOTA_AirBookRQ)

	if err != nil {
		return err
		//return err
	}

	fmt.Println(string(rev))

	var rs *RSBookingOTA_AirBookRS

	if errXML := xml.Unmarshal(rev, &rs); errXML != nil {
		return galaxylib.DefaultGalaxyError.FromError(1, errXML)
	}
	p.PNR = rs.RSBookingAirReservation.RSBookingBookingReferenceID.AttrID
	galaxylib.GalaxyLogger.Warningln(p.PNR)
	return nil
}

type BookingAirItinerary struct {
	BookingOriginDestinationOptions *BookingOriginDestinationOptions `xml:" OriginDestinationOptions,omitempty" json:"OriginDestinationOptions,omitempty"`
}

type BookingAirTraveler struct {
	AttrGender            string           `xml:" Gender,attr"  json:",omitempty"`
	AttrPassengerTypeCode string           `xml:" PassengerTypeCode,attr"  json:",omitempty"`
	BookingComment        *BookingComment  `xml:" Comment,omitempty" json:"Comment,omitempty"`
	BookingDocument       *BookingDocument `xml:" Document,omitempty" json:"Document,omitempty"`
	// BookingDocumentFlightBinding *BookingDocumentFlightBinding `xml:" DocumentFlightBinding,omitempty" json:"DocumentFlightBinding,omitempty"`
	BookingFlightSegmentRPHs *BookingFlightSegmentRPHs   `xml:" FlightSegmentRPHs,omitempty" json:"FlightSegmentRPHs,omitempty"`
	BookingPersonName        *BookingPersonName          `xml:" PersonName,omitempty" json:"PersonName,omitempty"`
	BookingTravelerRefNumber []*BookingTravelerRefNumber `xml:" TravelerRefNumber,omitempty" json:"TravelerRefNumber,omitempty"`
}

type BookingAirline struct {
	AttrCode string `xml:" Code,attr"  json:",omitempty"`
}

type BookingArrivalAirport struct {
	AttrLocationCode string `xml:" LocationCode,attr"  json:",omitempty"`
}

type BookingBookingClassAvail struct {
	AttrResBookDesigCode string `xml:" ResBookDesigCode,attr"  json:",omitempty"`
}

type BookingRoot struct {
	BookingOTA_AirBookRQ *BookingOTA_AirBookRQ `xml:" OTA_AirBookRQ,omitempty" json:"OTA_AirBookRQ,omitempty"`
}

type BookingComment struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type BookingContactInfo struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type BookingDepartureAirport struct {
	AttrLocationCode string `xml:" LocationCode,attr"  json:",omitempty"`
}

type BookingDocHolderFormattedName struct {
	BookingGivenName *BookingGivenName `xml:" GivenName,omitempty" json:"GivenName,omitempty"`
	BookingSurname   *BookingSurname   `xml:" Surname,omitempty" json:"Surname,omitempty"`
}

type BookingDocument struct {
	AttrBirthDate string `xml:" BirthDate,attr"  json:",omitempty"`
	// AttrDocHolderInd              string                         `xml:" DocHolderInd,attr"  json:",omitempty"` //是否为证件持有者
	AttrDocHolderNationality      string                         `xml:" DocHolderNationality,attr"  json:",omitempty"`
	AttrDocID                     string                         `xml:" DocID,attr"  json:",omitempty"`
	AttrDocIssueCountry           string                         `xml:" DocIssueCountry,attr"  json:",omitempty"`
	AttrDocType                   string                         `xml:" DocType,attr"  json:",omitempty"`
	AttrDocTypeDetail             string                         `xml:" DocTypeDetail,attr"  json:",omitempty"`
	AttrExpireDate                string                         `xml:" ExpireDate,attr"  json:",omitempty"`
	AttrGender                    string                         `xml:" Gender,attr"  json:",omitempty"`
	AttrRPH                       string                         `xml:" RPH,attr"  json:",omitempty"`
	BookingDocHolderFormattedName *BookingDocHolderFormattedName `xml:" DocHolderFormattedName,omitempty" json:"DocHolderFormattedName,omitempty"`
}

type BookingDocumentFlightBinding struct {
	BookingDocumentRPH      *BookingDocumentRPH      `xml:" DocumentRPH,omitempty" json:"DocumentRPH,omitempty"`
	BookingFlightSegmentRPH *BookingFlightSegmentRPH `xml:" FlightSegmentRPH,omitempty" json:"FlightSegmentRPH,omitempty"`
}

type BookingDocumentRPH struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type BookingEnvelopType struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type BookingEquipment struct {
	AttrAirEquipType string `xml:" AirEquipType,attr"  json:",omitempty"`
}

//航段信息
type BookingFlightSegment struct {
	AttrArrivalDateTime      string                    `xml:" ArrivalDateTime,attr"  json:",omitempty"`
	AttrCodeshareInd         string                    `xml:" CodeshareInd,attr"  json:",omitempty"` //是否共享航班，一般false
	AttrDepartureDateTime    string                    `xml:" DepartureDateTime,attr"  json:",omitempty"`
	AttrFlighttNumber        string                    `xml:" FlightNumber,attr"  json:",omitempty"`
	AttrRPH                  string                    `xml:" RPH,attr"  json:",omitempty"`         //航段编号,与旅客、ssr 等信息关联，多航段时，编号请勿重复。
	AttrSegmentType          string                    `xml:" SegmentType,attr"  json:",omitempty"` //NORMAL-普通航段 OPEN-不定期航段 ARRIVAL_UNKOWN_ARNK信息航段一般只写 NORMAL
	AttrStatus               string                    `xml:" Status,attr"  json:",omitempty"`      //即航段状态
	BookingArrivalAirport    *BookingArrivalAirport    `xml:" ArrivalAirport,omitempty" json:"ArrivalAirport,omitempty"`
	BookingBookingClassAvail *BookingBookingClassAvail `xml:" BookingClassAvail,omitempty" json:"BookingClassAvail,omitempty"`
	BookingDepartureAirport  *BookingDepartureAirport  `xml:" DepartureAirport,omitempty" json:"DepartureAirport,omitempty"`
	// BookingEquipment         *BookingEquipment         `xml:" Equipment,omitempty" json:"Equipment,omitempty"`
	BookingMarketingAirline *BookingMarketingAirline `xml:" MarketingAirline,omitempty" json:"MarketingAirline,omitempty"`
}

type BookingFlightSegmentRPH struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type BookingFlightSegmentRPHs struct {
	BookingFlightSegmentRPH *BookingFlightSegmentRPH `xml:" FlightSegmentRPH,omitempty" json:"FlightSegmentRPH,omitempty"`
}

type BookingGivenName struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type BookingMarketingAirline struct {
	AttrCode string `xml:" Code,attr"  json:",omitempty"`
}

type BookingOTA_AirBookRQ struct {
	BookingAirItinerary   *BookingAirItinerary   `xml:" AirItinerary,omitempty" json:"AirItinerary,omitempty"`
	BookingPOS            *BookingPOS            `xml:" POS,omitempty" json:"POS,omitempty"`
	BookingTPA_Extensions *BookingTPA_Extensions `xml:" TPA_Extensions,omitempty" json:"TPA_Extensions,omitempty"`
	BookingTicketing      *BookingTicketing      `xml:" Ticketing,omitempty" json:"Ticketing,omitempty"`
	BookingTravelerInfo   *BookingTravelerInfo   `xml:" TravelerInfo,omitempty" json:"TravelerInfo,omitempty"`
	XMLName               xml.Name               `xml:"OTA_AirBookRQ"`
}

type BookingOriginDestinationOption struct {
	BookingFlightSegment []*BookingFlightSegment `xml:" FlightSegment,omitempty" json:"FlightSegment,omitempty"`
}

type BookingOriginDestinationOptions struct {
	BookingOriginDestinationOption *BookingOriginDestinationOption `xml:" OriginDestinationOption,omitempty" json:"OriginDestinationOption,omitempty"`
}

type BookingOtherServiceInformation struct {
	AttrCode                 string                      `xml:" Code,attr"  json:",omitempty"`
	BookingAirline           *BookingAirline             `xml:" Airline,omitempty" json:"Airline,omitempty"`
	BookingText              *BookingText                `xml:" Text,omitempty" json:"Text,omitempty"`
	BookingTravelerRefNumber []*BookingTravelerRefNumber `xml:" TravelerRefNumber,omitempty" json:"TravelerRefNumber,omitempty"`
}

type BookingOtherServiceInformations struct {
	BookingOtherServiceInformation []*BookingOtherServiceInformation `xml:" OtherServiceInformation,omitempty" json:"OtherServiceInformation,omitempty"`
}

type BookingPOS struct {
	BookingSource *BookingSource `xml:" Source,omitempty" json:"Source,omitempty"`
}

type BookingPersonName struct {
	AttrLanguageType string          `xml:" LanguageType,attr"  json:",omitempty"`
	BookingSurname   *BookingSurname `xml:" Surname,omitempty" json:"Surname,omitempty"`
}

type BookingSource struct {
	AttrPseudoCityCode string `xml:" PseudoCityCode,attr"  json:",omitempty"`
}

type BookingSpecialReqDetails struct {
	BookingOtherServiceInformations *BookingOtherServiceInformations `xml:" OtherServiceInformations,omitempty" json:"OtherServiceInformations,omitempty"`
}

type BookingSurname struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type BookingTPA_Extensions struct {
	BookingContactInfo *BookingContactInfo `xml:" ContactInfo,omitempty" json:"ContactInfo,omitempty"`
	BookingEnvelopType *BookingEnvelopType `xml:" EnvelopType,omitempty" json:"EnvelopType,omitempty"`
}

type BookingText struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type BookingTicketing struct {
	AttrTicketTimeLimit string `xml:" TicketTimeLimit,attr"  json:",omitempty"`
}

type BookingTravelerInfo struct {
	BookingAirTraveler       []*BookingAirTraveler     `xml:" AirTraveler,omitempty" json:"AirTraveler,omitempty"`
	BookingSpecialReqDetails *BookingSpecialReqDetails `xml:" SpecialReqDetails,omitempty" json:"SpecialReqDetails,omitempty"`
}

type BookingTravelerRefNumber struct {
	AttrRPH string `xml:" RPH,attr"  json:",omitempty"` //旅客编号，请不要输入重复的值若分 pnr 预订，请必填此项
}

/*************************************
	Response type
********************************************/

type RSBookingAirItinerary struct {
	RSBookingFlightSegments *RSBookingFlightSegments `xml:" FlightSegments,omitempty" json:"FlightSegments,omitempty"`
}

type RSBookingAirReservation struct {
	RSBookingAirItinerary       *RSBookingAirItinerary       `xml:" AirItinerary,omitempty" json:"AirItinerary,omitempty"`
	RSBookingBookingReferenceID *RSBookingBookingReferenceID `xml:" BookingReferenceID,omitempty" json:"BookingReferenceID,omitempty"`
	RSBookingComment            []*RSBookingComment          `xml:" Comment,omitempty" json:"Comment,omitempty"`
}

type RSBookingArrivalAirport struct {
	AttrLocationCode string `xml:" LocationCode,attr"  json:",omitempty"`
}

type RSBookingBookingClassAvail struct {
	AttrResBookDesigCode string `xml:" ResBookDesigCode,attr"  json:",omitempty"`
}

type RSBookingBookingReferenceID struct {
	AttrID         string `xml:" ID,attr"  json:",omitempty"`
	AttrID_Context string `xml:" ID_Context,attr"  json:",omitempty"`
}

type RSBookingRoot struct {
	RSBookingOTA_AirBookRS *RSBookingOTA_AirBookRS `xml:" OTA_AirBookRS,omitempty" json:"OTA_AirBookRS,omitempty"`
}

type RSBookingComment struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type RSBookingDepartureAirport struct {
	AttrLocationCode string `xml:" LocationCode,attr"  json:",omitempty"`
}

type RSBookingFlightSegment struct {
	AttrArrivalDateTime        string                      `xml:" ArrivalDateTime,attr"  json:",omitempty"`
	AttrCodeshareInd           string                      `xml:" CodeshareInd,attr"  json:",omitempty"`
	AttrDepartureDateTime      string                      `xml:" DepartureDateTime,attr"  json:",omitempty"`
	AttrFlightNumber           string                      `xml:" FlightNumber,attr"  json:",omitempty"`
	AttrNumberInParty          string                      `xml:" NumberInParty,attr"  json:",omitempty"`
	AttrSegmentType            string                      `xml:" SegmentType,attr"  json:",omitempty"`
	AttrStatus                 string                      `xml:" Status,attr"  json:",omitempty"`
	RSBookingArrivalAirport    *RSBookingArrivalAirport    `xml:" ArrivalAirport,omitempty" json:"ArrivalAirport,omitempty"`
	RSBookingBookingClassAvail *RSBookingBookingClassAvail `xml:" BookingClassAvail,omitempty" json:"BookingClassAvail,omitempty"`
	RSBookingDepartureAirport  *RSBookingDepartureAirport  `xml:" DepartureAirport,omitempty" json:"DepartureAirport,omitempty"`
	RSBookingMarketingAirline  *RSBookingMarketingAirline  `xml:" MarketingAirline,omitempty" json:"MarketingAirline,omitempty"`
	RSBookingOperatingAirline  *RSBookingOperatingAirline  `xml:" OperatingAirline,omitempty" json:"OperatingAirline,omitempty"`
}

type RSBookingFlightSegments struct {
	RSBookingFlightSegment *RSBookingFlightSegment `xml:" FlightSegment,omitempty" json:"FlightSegment,omitempty"`
}

type RSBookingMarketingAirline struct {
	AttrCode string `xml:" Code,attr"  json:",omitempty"`
}

type RSBookingOTA_AirBookRS struct {
	RSBookingAirReservation *RSBookingAirReservation `xml:" AirReservation,omitempty" json:"AirReservation,omitempty"`
}

type RSBookingOperatingAirline struct {
}
