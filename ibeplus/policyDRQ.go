package ibeplus

import (
	"encoding/xml"
	"fmt"
	"ibepluslib/models"
	"strconv"
	"time"

	"github.com/otwdev/galaxylib"
)

type PolicyDRQ struct {
	Order *models.OrderInfo
}

func NewPolicyDRQ(order *models.OrderInfo) *PolicyDRQ {
	return &PolicyDRQ{
		Order: order,
	}
}

func (p *PolicyDRQ) PolicyD() *PolicyDRSRoot {

	pnr := p.Order.PnrInofs[0]

	rq := &PolicyDRQFareInterface{}
	rq.PolicyDRQInput = &PolicyDRQInput{}
	rq.PolicyDRQInput.PolicyDRQHeaderIn = &PolicyDRQHeaderIn{}
	rq.PolicyDRQInput.PolicyDRQRequest = &PolicyDRQRequest{}

	rq.PolicyDRQInput.PolicyDRQHeaderIn.PolicyDRQSysCode = &PolicyDRQSysCode{
		Text: "CRS",
	}

	rq.PolicyDRQInput.PolicyDRQHeaderIn.PolicyDRQChannelID = &PolicyDRQChannelID{
		Text: "1E",
	}

	rq.PolicyDRQInput.PolicyDRQHeaderIn.PolicyDRQCommandType = &PolicyDRQCommandType{
		Text: "FLS",
	}

	rq.PolicyDRQInput.PolicyDRQHeaderIn.PolicyDRQChannelType = &PolicyDRQChannelType{
		Text: "COMMON",
	}

	rq.PolicyDRQInput.PolicyDRQHeaderIn.PolicyDRQLanguage = &PolicyDRQLanguage{
		Text: "CN",
	}

	agency := &PolicyDRQAgency{}

	agency.PolicyDRQOfficeId = &PolicyDRQOfficeId{
		Text: pnr.OfficeNumber,
	}

	agency.PolicyDRQPid = &PolicyDRQPid{
		Text: galaxylib.GalaxyCfgFile.MustValue("ibe", "defaultPid"), //"394541",
	}

	agency.PolicyDRQCity = &PolicyDRQCity{
		Text: pnr.AgencyCity,
	}

	rq.PolicyDRQInput.PolicyDRQHeaderIn.PolicyDRQAgency = agency

	rq.PolicyDRQInput.PolicyDRQHeaderIn.PolicyDRQDistributor = &PolicyDRQDistributor{
		Text: "运营二部/同业线下",
	}

	rq.PolicyDRQInput.PolicyDRQRequest.PolicyDRQFlightShopRequest = &PolicyDRQFlightShopRequest{}

	for _, v := range pnr.FlightSegments {
		od := &PolicyDRQOriginDestinationInfo{}
		od.PolicyDRQOri = &PolicyDRQOri{
			Text: v.DepartCityCode,
		}
		od.PolicyDRQDes = &PolicyDRQDes{
			Text: v.ArriveCityCode,
		}
		// d, t := formatRQTime(v)
		od.PolicyDRQDepartureDate = &PolicyDRQDepartureDate{
			Text: v.DepDate02Jan06(),
		}
		od.PolicyDRQDepartureTime1 = &PolicyDRQDepartureTime1{
			Text: v.DepTime1504(),
		}
		od.PolicyDRQCarrier = &PolicyDRQCarrier{
			Text: v.MarketingAirLine,
		}
		rq.PolicyDRQInput.PolicyDRQRequest.PolicyDRQFlightShopRequest.PolicyDRQOriginDestinationInfo =
			append(rq.PolicyDRQInput.PolicyDRQRequest.PolicyDRQFlightShopRequest.PolicyDRQOriginDestinationInfo, od)

		//rq.PolicyDRQInput.PolicyDRQRequest.PolicyDRQFlightShopRequest.PolicyDRQTravelPreferences
	}

	travlePre := &PolicyDRQTravelPreferences{}
	travlePre.PolicyDRQIsDirectFlightOnly = &PolicyDRQIsDirectFlightOnly{
		Text: "true",
	}
	travlePre.PolicyDRQDisplayCurrCode = &PolicyDRQDisplayCurrCode{
		Text: "CNY",
	}
	travlePre.PolicyDRQCurrCode = &PolicyDRQCurrCode{
		Text: "CNY",
	}
	travelType := "OW"
	if len(pnr.FlightSegments) > 1 {
		travelType = "RT"
	}
	travlePre.PolicyDRQJourneyType = &PolicyDRQJourneyType{
		Text: travelType,
	}

	adt, chd := pnr.PersonQuantity()

	travlePre.passengerPackage(adt, models.Adult)
	travlePre.passengerPackage(chd, models.Child)

	rq.PolicyDRQInput.PolicyDRQRequest.PolicyDRQFlightShopRequest.PolicyDRQTravelPreferences = travlePre

	//journey := p.makeJourney(pnr)

	//rq.PolicyDRQInput.PolicyDRQRequest.PolicyDRQFlightShopRequest.PolicyDRQAvJourneys = journey

	rq.PolicyDRQInput.PolicyDRQRequest.PolicyDRQFlightShopRequest.PolicyDRQOption = &PolicyDRQOption{}
	// rq.PolicyDRQInput.PolicyDRQRequest.PolicyDRQFlightShopRequest.PolicyDRQOption.PolicyDRQLowestOrAll = &PolicyDRQLowestOrAll{
	// 	Text: "A",
	// }

	return p.reqeust(rq)

}

func (p *PolicyDRQ) makeJourney(pnr *models.PnrInfo) *PolicyDRQAvJourneys {

	//单程
	flt := pnr.FlightSegments[0]

	journey := &PolicyDRQAvJourneys{}
	journey.PolicyDRQAvJourney = &PolicyDRQAvJourney{}
	journey.PolicyDRQAvJourney.PolicyDRQArr = &PolicyDRQArr{
		Text: flt.ArriveCityCode,
	}
	journey.PolicyDRQAvJourney.PolicyDRQDep = &PolicyDRQDep{
		Text: flt.DepartCityCode,
	}
	avOpt := &PolicyDRQAvOpt{}

	avOpt.PolicyDRQFlt = &PolicyDRQFlt{}
	avOpt.PolicyDRQFlt.PolicyDRQAirline = &PolicyDRQAirline{
		Text: flt.AirlineCode(),
	}
	avOpt.PolicyDRQFlt.PolicyDRQArr = &PolicyDRQArr{
		Text: flt.ArriveCityCode,
	}
	avOpt.PolicyDRQFlt.PolicyDRQDep = &PolicyDRQDep{
		Text: flt.DepartCityCode,
	}

	// dt, dm := formatRQTime(flt)

	avOpt.PolicyDRQFlt.PolicyDRQDt = &PolicyDRQDt{
		Text: flt.DepDate02Jan06(),
	}
	avOpt.PolicyDRQFlt.PolicyDRQDeptm = &PolicyDRQDeptm{
		Text: flt.DepTime1504(),
	}

	t, _ := time.Parse("15:04", flt.ArrTime)
	retT := t.Format("1504")
	avOpt.PolicyDRQFlt.PolicyDRQArrtm = &PolicyDRQArrtm{
		Text: retT,
	}
	avOpt.PolicyDRQFlt.PolicyDRQFltNo = &PolicyDRQFltNo{
		Text: flt.AirlineNum(),
	}

	journey.PolicyDRQAvJourney.PolicyDRQAvOpt = append(journey.PolicyDRQAvJourney.PolicyDRQAvOpt, avOpt)

	return journey
}

const policyDURL = "http://ibeplus.travelsky.com/ota/xml/AirFlightShopPolicy/D"

func (t *PolicyDRQ) reqeust(rq *PolicyDRQFareInterface) *PolicyDRSRoot {
	ibe := NewIBE(policyDURL, rq)
	buf, err := ibe.Reqeust()
	if err != nil {
		galaxylib.GalaxyLogger.Errorln(err)
	}
	root := &PolicyDRSRoot{}
	// errJSON := json.Unmarshal(buf, &root.PolicyDRSFareInterface)
	errJSON := xml.Unmarshal(buf, &root.PolicyDRSFareInterface)
	if errJSON != nil {
		fmt.Println(errJSON)
	}
	//fmt.Println(string(buf))
	return root
}

func (t *PolicyDRQTravelPreferences) passengerPackage(quantity int, ty models.PersonType) {
	if quantity == 0 {
		return
	}
	psn := &PolicyDRQPassenger{}
	psn.PolicyDRQNumber = &PolicyDRQNumber{
		Text: strconv.Itoa(quantity),
	}
	psn.PolicyDRQType = &PolicyDRQType{
		Text: string(ty),
	}
	t.PolicyDRQPassenger = append(t.PolicyDRQPassenger, psn)
}

// func formatRQTime(f *models.FlightSegment) (string, string) {
// 	d, _ := time.Parse("2006-01-02", f.FlyDate)
// 	revD := d.Format("02Jan06")
// 	revD = strings.ToUpper(revD)
// 	t, _ := time.Parse("15:04", f.DepTime)
// 	retT := t.Format("1504")
// 	return revD, retT
// }

type PolicyDRQASR struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQAgency struct {
	PolicyDRQCity     *PolicyDRQCity     `xml:" city,omitempty" json:"city,omitempty"`
	PolicyDRQOfficeId *PolicyDRQOfficeId `xml:" officeId,omitempty" json:"officeId,omitempty"`
	PolicyDRQPid      *PolicyDRQPid      `xml:" pid,omitempty" json:"pid,omitempty"`
}

type PolicyDRQAvJourney struct {
	PolicyDRQArr   *PolicyDRQArr     `xml:" arr,omitempty" json:"arr,omitempty"`
	PolicyDRQAvOpt []*PolicyDRQAvOpt `xml:" AvOpt,omitempty" json:"AvOpt,omitempty"`
	PolicyDRQDep   *PolicyDRQDep     `xml:" dep,omitempty" json:"dep,omitempty"`
	PolicyDRQDt    *PolicyDRQDt      `xml:" dt,omitempty" json:"dt,omitempty"`
	PolicyDRQRPH   *PolicyDRQRPH     `xml:" RPH,omitempty" json:"RPH,omitempty"`
	PolicyDRQWeek  *PolicyDRQWeek    `xml:" week,omitempty" json:"week,omitempty"`
}

type PolicyDRQAvJourneys struct {
	PolicyDRQAvJourney *PolicyDRQAvJourney `xml:" AvJourney,omitempty" json:"AvJourney,omitempty"`
	PolicyDRQOffice    *PolicyDRQOffice    `xml:" office,omitempty" json:"office,omitempty"`
	PolicyDRQRPH       *PolicyDRQRPH       `xml:" RPH,omitempty" json:"RPH,omitempty"`
}

type PolicyDRQAvOpt struct {
	PolicyDRQFlt *PolicyDRQFlt `xml:" Flt,omitempty" json:"Flt,omitempty"`
	PolicyDRQRPH *PolicyDRQRPH `xml:" RPH,omitempty" json:"RPH,omitempty"`
}

// type PolicyDRQChidleyRoot314159 struct {
// 	PolicyDRQFareInterface *PolicyDRQFareInterface `xml:" FareInterface,omitempty" json:"FareInterface,omitempty"`
// }

type PolicyDRQDepartureDate struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQDepartureTime1 struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQDistributor struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQFareInterface struct {
	PolicyDRQInput *PolicyDRQInput `xml:" Input,omitempty" json:"Input,omitempty"`
	XMLName        xml.Name        `xml:"FareInterface"`
}

// type PolicyDRQFareInterface struct {
// 	PolicyDRQInput *PolicyDRQInput `xml:" Input,omitempty" json:"Input,omitempty"`
// }

type PolicyDRQFlightShopRequest struct {
	PolicyDRQAvJourneys            *PolicyDRQAvJourneys              `xml:" AvJourneys,omitempty" json:"AvJourneys,omitempty"`
	PolicyDRQOption                *PolicyDRQOption                  `xml:" Option,omitempty" json:"Option,omitempty"`
	PolicyDRQOriginDestinationInfo []*PolicyDRQOriginDestinationInfo `xml:" OriginDestinationInfo,omitempty" json:"OriginDestinationInfo,omitempty"`
	PolicyDRQTravelPreferences     *PolicyDRQTravelPreferences       `xml:" TravelPreferences,omitempty" json:"TravelPreferences,omitempty"`
}

type PolicyDRQFlt struct {
	PolicyDRQASR       *PolicyDRQASR       `xml:" ASR,omitempty" json:"ASR,omitempty"`
	PolicyDRQAirline   *PolicyDRQAirline   `xml:" airline,omitempty" json:"airline,omitempty"`
	PolicyDRQArr       *PolicyDRQArr       `xml:" arr,omitempty" json:"arr,omitempty"`
	PolicyDRQArrtm     *PolicyDRQArrtm     `xml:" arrtm,omitempty" json:"arrtm,omitempty"`
	PolicyDRQClass     []*PolicyDRQClass   `xml:" class,omitempty" json:"class,omitempty"`
	PolicyDRQCodeshare *PolicyDRQCodeshare `xml:" codeshare,omitempty" json:"codeshare,omitempty"`
	PolicyDRQDep       *PolicyDRQDep       `xml:" dep,omitempty" json:"dep,omitempty"`
	PolicyDRQDeptm     *PolicyDRQDeptm     `xml:" deptm,omitempty" json:"deptm,omitempty"`
	PolicyDRQDev       *PolicyDRQDev       `xml:" dev,omitempty" json:"dev,omitempty"`
	PolicyDRQDt        *PolicyDRQDt        `xml:" dt,omitempty" json:"dt,omitempty"`
	PolicyDRQEt        *PolicyDRQEt        `xml:" et,omitempty" json:"et,omitempty"`
	PolicyDRQFltNo     *PolicyDRQFltNo     `xml:" fltNo,omitempty" json:"fltNo,omitempty"`
	PolicyDRQLnk       *PolicyDRQLnk       `xml:" lnk,omitempty" json:"lnk,omitempty"`
	PolicyDRQMeal      *PolicyDRQMeal      `xml:" meal,omitempty" json:"meal,omitempty"`
	PolicyDRQPgind     *PolicyDRQPgind     `xml:" pgind,omitempty" json:"pgind,omitempty"`
	PolicyDRQRPH       *PolicyDRQRPH       `xml:" RPH,omitempty" json:"RPH,omitempty"`
	PolicyDRQRoutno    *PolicyDRQRoutno    `xml:" routno,omitempty" json:"routno,omitempty"`
	PolicyDRQStop      *PolicyDRQStop      `xml:" stop,omitempty" json:"stop,omitempty"`
	PolicyDRQSubid     *PolicyDRQSubid     `xml:" subid,omitempty" json:"subid,omitempty"`
	PolicyDRQTerm      *PolicyDRQTerm      `xml:" term,omitempty" json:"term,omitempty"`
	PolicyDRQTpm       *PolicyDRQTpm       `xml:" tpm,omitempty" json:"tpm,omitempty"`
	PolicyDRQWeek      *PolicyDRQWeek      `xml:" week,omitempty" json:"week,omitempty"`
}

type PolicyDRQGroupType struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQHeaderIn struct {
	PolicyDRQAgency      *PolicyDRQAgency      `xml:" Agency,omitempty" json:"Agency,omitempty"`
	PolicyDRQChannelID   *PolicyDRQChannelID   `xml:" channelID,omitempty" json:"channelID,omitempty"`
	PolicyDRQChannelType *PolicyDRQChannelType `xml:" channelType,omitempty" json:"channelType,omitempty"`
	PolicyDRQCommandType *PolicyDRQCommandType `xml:" commandType,omitempty" json:"commandType,omitempty"`
	PolicyDRQDistributor *PolicyDRQDistributor `xml:" Distributor,omitempty" json:"Distributor,omitempty"`
	PolicyDRQGroupType   *PolicyDRQGroupType   `xml:" GroupType,omitempty" json:"GroupType,omitempty"`
	PolicyDRQLanguage    *PolicyDRQLanguage    `xml:" language,omitempty" json:"language,omitempty"`
	PolicyDRQSysCode     *PolicyDRQSysCode     `xml:" sysCode,omitempty" json:"sysCode,omitempty"`
}

type PolicyDRQInput struct {
	PolicyDRQHeaderIn *PolicyDRQHeaderIn `xml:" HeaderIn,omitempty" json:"HeaderIn,omitempty"`
	PolicyDRQRequest  *PolicyDRQRequest  `xml:" Request,omitempty" json:"Request,omitempty"`
}

type PolicyDRQOption struct {
	PolicyDRQFareSource        *PolicyDRQFareSource        `xml:" fareSource,omitempty" json:"fareSource,omitempty"`
	PolicyDRQFcFeature         *PolicyDRQFcFeature         `xml:" fcFeature,omitempty" json:"fcFeature,omitempty"`
	PolicyDRQFormat            *PolicyDRQFormat            `xml:" format,omitempty" json:"format,omitempty"`
	PolicyDRQIsAvNeeded        *PolicyDRQIsAvNeeded        `xml:" isAvNeeded,omitempty" json:"isAvNeeded,omitempty"`
	PolicyDRQIsFaresNeeded     *PolicyDRQIsFaresNeeded     `xml:" isFaresNeeded,omitempty" json:"isFaresNeeded,omitempty"`
	PolicyDRQIsPSnNeeded       *PolicyDRQIsPSnNeeded       `xml:" isPSnNeeded,omitempty" json:"isPSnNeeded,omitempty"`
	PolicyDRQIsPsAvBindsNeeded *PolicyDRQIsPsAvBindsNeeded `xml:" isPsAvBindsNeeded,omitempty" json:"isPsAvBindsNeeded,omitempty"`
	PolicyDRQLowestOrAll       *PolicyDRQLowestOrAll       `xml:" lowestOrAll,omitempty" json:"lowestOrAll,omitempty"`
	PolicyDRQRuleTypeNeeded    *PolicyDRQRuleTypeNeeded    `xml:" ruleTypeNeeded,omitempty" json:"ruleTypeNeeded,omitempty"`
}

type PolicyDRQOriginDestinationInfo struct {
	PolicyDRQCarrier        *PolicyDRQCarrier        `xml:" carrier,omitempty" json:"carrier,omitempty"`
	PolicyDRQDepartureDate  *PolicyDRQDepartureDate  `xml:" DepartureDate,omitempty" json:"DepartureDate,omitempty"`
	PolicyDRQDepartureTime1 *PolicyDRQDepartureTime1 `xml:" DepartureTime1,omitempty" json:"DepartureTime1,omitempty"`
	PolicyDRQDes            *PolicyDRQDes            `xml:" des,omitempty" json:"des,omitempty"`
	PolicyDRQOri            *PolicyDRQOri            `xml:" ori,omitempty" json:"ori,omitempty"`
}

type PolicyDRQRPH struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQRequest struct {
	PolicyDRQFlightShopRequest *PolicyDRQFlightShopRequest `xml:" FlightShopRequest,omitempty" json:"FlightShopRequest,omitempty"`
}

type PolicyDRQTravelPreferences struct {
	PolicyDRQCabinClass         *PolicyDRQCabinClass         `xml:" cabinClass,omitempty" json:"cabinClass,omitempty"`
	PolicyDRQCurrCode           *PolicyDRQCurrCode           `xml:" currCode,omitempty" json:"currCode,omitempty"`
	PolicyDRQDisplayCurrCode    *PolicyDRQDisplayCurrCode    `xml:" displayCurrCode,omitempty" json:"displayCurrCode,omitempty"`
	PolicyDRQIsDealModel        *PolicyDRQIsDealModel        `xml:" isDealModel,omitempty" json:"isDealModel,omitempty"`
	PolicyDRQIsDirectFlightOnly *PolicyDRQIsDirectFlightOnly `xml:" isDirectFlightOnly,omitempty" json:"isDirectFlightOnly,omitempty"`
	PolicyDRQJourneyType        *PolicyDRQJourneyType        `xml:" journeyType,omitempty" json:"journeyType,omitempty"`
	PolicyDRQPassenger          []*PolicyDRQPassenger        `xml:" passenger,omitempty" json:"passenger,omitempty"`
}

type PolicyDRQAirline struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQArr struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQArrtm struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQAv struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQCabinClass struct {
}

type PolicyDRQCarrier struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQChannelID struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQChannelType struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQCity struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQClass struct {
	PolicyDRQAv   *PolicyDRQAv   `xml:" av,omitempty" json:"av,omitempty"`
	PolicyDRQName *PolicyDRQName `xml:" name,omitempty" json:"name,omitempty"`
}

type PolicyDRQCodeshare struct {
	PolicyDRQAirline *PolicyDRQAirline `xml:" airline,omitempty" json:"airline,omitempty"`
	PolicyDRQFltno   *PolicyDRQFltno   `xml:" fltno,omitempty" json:"fltno,omitempty"`
}

type PolicyDRQCommandType struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQCurrCode struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQDep struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQDeptm struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQDes struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQDev struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQDisplayCurrCode struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQDt struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQEt struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQFareSource struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQFcFeature struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQFltNo struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQFltno struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQFormat struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQIsAvNeeded struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQIsDealModel struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQIsDirectFlightOnly struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQIsFaresNeeded struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQIsPSnNeeded struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQIsPsAvBindsNeeded struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQJourneyType struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQLanguage struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQLnk struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQLowestOrAll struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQMeal struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQName struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQNumber struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQOffice struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQOfficeId struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQOri struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQPassenger struct {
	PolicyDRQNumber *PolicyDRQNumber `xml:" number,omitempty" json:"number,omitempty"`
	PolicyDRQType   *PolicyDRQType   `xml:" type,omitempty" json:"type,omitempty"`
}

type PolicyDRQPgind struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQPid struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQRoutno struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQRuleTypeNeeded struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQStop struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQSubid struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQSysCode struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQTerm struct {
	PolicyDRQArr *PolicyDRQArr `xml:" arr,omitempty" json:"arr,omitempty"`
	PolicyDRQDep *PolicyDRQDep `xml:" dep,omitempty" json:"dep,omitempty"`
}

type PolicyDRQTpm struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQType struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRQWeek struct {
	Text string `xml:",chardata" json:",omitempty"`
}
