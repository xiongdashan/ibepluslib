package ibeplus

import (
	"ibepluslib/models"

	"github.com/otwdev/galaxylib"
)

func NewPolicyDRS() *PolicyDRSRoot {
	return &PolicyDRSRoot{}
}

func (p *PolicyDRSRoot) GetPNRPrice(pnr *models.PnrInfo) *galaxylib.GalaxyError {

	err := p.PolicyDRSFareInterface.PolicyDRSOutput.PolicyDRSResult.PolicyDRSError

	if err != nil {
		return galaxylib.DefaultGalaxyError.FromText(1, err.PolicyDRSMessage.Text)
	}

	psAry := p.PolicyDRSFareInterface.PolicyDRSOutput.PolicyDRSResult.PolicyDRSFlightShopResult.PolicyDRSPSn.PolicyDRSPS

	flt := pnr.FlightSegments[0]

	for _, v := range psAry {
		for _, r := range v.PolicyDRSRouts {
			route := r.PolicyDRSRout
			depd := flt.DepDate02Jan2006()
			arrd := flt.ArrDate02Jan2006()
			if route.PolicyDRSFltNo.Text == flt.FlyNo &&
				route.PolicyDRSCarr.Text == flt.MarketingAirLine &&
				route.PolicyDRSDepartureAirport.Text == flt.DepartCityCode &&
				route.PolicyDRSArrivalAirport.Text == flt.ArriveCityCode &&
				route.PolicyDRSDepartureDate.Text == depd &&
				route.PolicyDRSDepartureTime.Text == flt.DepTime1504() &&
				route.PolicyDRSArrivalDate.Text == arrd &&
				route.PolicyDRSArrivalTime.Text == flt.ArrTime1504() {
				//return v, r.PolicyDRSRout, nil
				pnr.FlightSegments[0].Cabin = route.PolicyDRSBkClass.Text
				pnr = p.fillingPNR(v, r.PolicyDRSRout, pnr)
				return nil
			}
		}
	}
	return galaxylib.DefaultGalaxyError.FromText(1, "末找到数据...")

}

func (p *PolicyDRSRoot) fillingPNR(ps *PolicyDRSPS, route *PolicyDRSRout, pnr *models.PnrInfo) *models.PnrInfo {
	//一个成人价格
	price := &models.PnrPrice{}
	price.Type = string(models.Adult)
	price.UpPrice = galaxylib.DefaultGalaxyConverter.MustFloat(ps.PolicyDRSDisAmt.Text)
	for _, fax := range ps.PolicyDRSTaxes.PolicyDRSTax {
		if fax.PolicyDRSCode.Text == "CN" {
			price.UpFax = galaxylib.DefaultGalaxyConverter.MustFloat(fax.PolicyDRSAmt.Text)
			continue
		}
		if fax.PolicyDRSCode.Text == "YQ" {
			price.YQ = galaxylib.DefaultGalaxyConverter.MustFloat(fax.PolicyDRSAmt.Text)
		}

	}
	pnr.Price = append(pnr.Price, price)
	return pnr
}

type PolicyDRSRoot struct {
	PolicyDRSFareInterface *PolicyDRSFareInterface `xml:" FareInterface,omitempty" json:"FareInterface,omitempty"`
}

type PolicyDRSA struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSASR struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSAvJourney struct {
	PolicyDRSArr   *PolicyDRSArr     `xml:" arr,omitempty" json:"arr,omitempty"`
	PolicyDRSAvOpt []*PolicyDRSAvOpt `xml:" AvOpt,omitempty" json:"AvOpt,omitempty"`
	PolicyDRSDep   *PolicyDRSDep     `xml:" dep,omitempty" json:"dep,omitempty"`
	PolicyDRSDt    *PolicyDRSDt      `xml:" dt,omitempty" json:"dt,omitempty"`
	PolicyDRSRPH   *PolicyDRSRPH     `xml:" RPH,omitempty" json:"RPH,omitempty"`
	PolicyDRSWeek  *PolicyDRSWeek    `xml:" week,omitempty" json:"week,omitempty"`
}

type PolicyDRSAvJourneys struct {
	PolicyDRSAvJourney *PolicyDRSAvJourney `xml:" AvJourney,omitempty" json:"AvJourney,omitempty"`
	PolicyDRSOffice    *PolicyDRSOffice    `xml:" office,omitempty" json:"office,omitempty"`
	PolicyDRSRPH       *PolicyDRSRPH       `xml:" RPH,omitempty" json:"RPH,omitempty"`
}

type PolicyDRSAvOpt struct {
	PolicyDRSFlt *PolicyDRSFlt `xml:" Flt,omitempty" json:"Flt,omitempty"`
	PolicyDRSRPH *PolicyDRSRPH `xml:" RPH,omitempty" json:"RPH,omitempty"`
}

type PolicyDRSB struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSCombineRule struct {
	PolicyDRSCombineOpenJawFlag   *PolicyDRSCombineOpenJawFlag   `xml:" combineOpenJawFlag,omitempty" json:"combineOpenJawFlag,omitempty"`
	PolicyDRSCombineRoundTripFlag *PolicyDRSCombineRoundTripFlag `xml:" combineRoundTripFlag,omitempty" json:"combineRoundTripFlag,omitempty"`
	PolicyDRSFareCombinationFlag  *PolicyDRSFareCombinationFlag  `xml:" fareCombinationFlag,omitempty" json:"fareCombinationFlag,omitempty"`
}

type PolicyDRSContent struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSEI struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSError struct {
	PolicyDRSCode    *PolicyDRSCode    `xml:" code,omitempty" json:"code,omitempty"`
	PolicyDRSMessage *PolicyDRSMessage `xml:" message,omitempty" json:"message,omitempty"`
}

type PolicyDRSFC struct {
	PolicyDRSDes         *PolicyDRSDes         `xml:" des,omitempty" json:"des,omitempty"`
	PolicyDRSDisAmt      *PolicyDRSDisAmt      `xml:" disAmt,omitempty" json:"disAmt,omitempty"`
	PolicyDRSDisCurrCode *PolicyDRSDisCurrCode `xml:" disCurrCode,omitempty" json:"disCurrCode,omitempty"`
	PolicyDRSEndorsement *PolicyDRSEndorsement `xml:" endorsement,omitempty" json:"endorsement,omitempty"`
	PolicyDRSFareBasis   *PolicyDRSFareBasis   `xml:" fareBasis,omitempty" json:"fareBasis,omitempty"`
	PolicyDRSFareBind    *PolicyDRSFareBind    `xml:" FareBind,omitempty" json:"FareBind,omitempty"`
	PolicyDRSOri         *PolicyDRSOri         `xml:" ori,omitempty" json:"ori,omitempty"`
	PolicyDRSOriAmt      *PolicyDRSOriAmt      `xml:" oriAmt,omitempty" json:"oriAmt,omitempty"`
	PolicyDRSSecInfo     *PolicyDRSSecInfo     `xml:" SecInfo,omitempty" json:"SecInfo,omitempty"`
	PolicyDRSYFares      *PolicyDRSYFares      `xml:" YFares,omitempty" json:"YFares,omitempty"`
}

type PolicyDRSFCs struct {
	PolicyDRSFC *PolicyDRSFC `xml:" FC,omitempty" json:"FC,omitempty"`
}

type PolicyDRSFareBind struct {
	PolicyDRSFareRPH *PolicyDRSFareRPH `xml:" fareRPH,omitempty" json:"fareRPH,omitempty"`
	PolicyDRSSysType *PolicyDRSSysType `xml:" sysType,omitempty" json:"sysType,omitempty"`
}

type PolicyDRSFareInterface struct {
	PolicyDRSOutput *PolicyDRSOutput `xml:" Output,omitempty" json:"Output,omitempty"`
}

type PolicyDRSFbrDtls struct {
}

type PolicyDRSFc struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSFlightShopResult struct {
	PolicyDRSAvJourneys     *PolicyDRSAvJourneys     `xml:" AvJourneys,omitempty" json:"AvJourneys,omitempty"`
	PolicyDRSFbrDtls        *PolicyDRSFbrDtls        `xml:" FbrDtls,omitempty" json:"FbrDtls,omitempty"`
	PolicyDRSNFares         *PolicyDRSNFares         `xml:" NFares,omitempty" json:"NFares,omitempty"`
	PolicyDRSPFares         *PolicyDRSPFares         `xml:" PFares,omitempty" json:"PFares,omitempty"`
	PolicyDRSPSn            *PolicyDRSPSn            `xml:" PSn,omitempty" json:"PSn,omitempty"`
	PolicyDRSPolicyBindings *PolicyDRSPolicyBindings `xml:" PolicyBindings,omitempty" json:"PolicyBindings,omitempty"`
	PolicyDRSPsAvBinds      *PolicyDRSPsAvBinds      `xml:" PsAvBinds,omitempty" json:"PsAvBinds,omitempty"`
	PolicyDRSRefundRules    *PolicyDRSRefundRules    `xml:" RefundRules,omitempty" json:"RefundRules,omitempty"`
	PolicyDRSReissueRules   *PolicyDRSReissueRules   `xml:" ReissueRules,omitempty" json:"ReissueRules,omitempty"`
	PolicyDRSRules          *PolicyDRSRules          `xml:" Rules,omitempty" json:"Rules,omitempty"`
	PolicyDRSShopWarning    *PolicyDRSShopWarning    `xml:" ShopWarning,omitempty" json:"ShopWarning,omitempty"`
}

type PolicyDRSFlt struct {
	PolicyDRSASR       *PolicyDRSASR       `xml:" ASR,omitempty" json:"ASR,omitempty"`
	PolicyDRSAirline   *PolicyDRSAirline   `xml:" airline,omitempty" json:"airline,omitempty"`
	PolicyDRSArr       *PolicyDRSArr       `xml:" arr,omitempty" json:"arr,omitempty"`
	PolicyDRSArrtm     *PolicyDRSArrtm     `xml:" arrtm,omitempty" json:"arrtm,omitempty"`
	PolicyDRSClass     []*PolicyDRSClass   `xml:" class,omitempty" json:"class,omitempty"`
	PolicyDRSCodeshare *PolicyDRSCodeshare `xml:" codeshare,omitempty" json:"codeshare,omitempty"`
	PolicyDRSDep       *PolicyDRSDep       `xml:" dep,omitempty" json:"dep,omitempty"`
	PolicyDRSDeptm     *PolicyDRSDeptm     `xml:" deptm,omitempty" json:"deptm,omitempty"`
	PolicyDRSDev       *PolicyDRSDev       `xml:" dev,omitempty" json:"dev,omitempty"`
	PolicyDRSDt        *PolicyDRSDt        `xml:" dt,omitempty" json:"dt,omitempty"`
	PolicyDRSEt        *PolicyDRSEt        `xml:" et,omitempty" json:"et,omitempty"`
	PolicyDRSFltNo     *PolicyDRSFltNo     `xml:" fltNo,omitempty" json:"fltNo,omitempty"`
	PolicyDRSLnk       *PolicyDRSLnk       `xml:" lnk,omitempty" json:"lnk,omitempty"`
	PolicyDRSMeal      *PolicyDRSMeal      `xml:" meal,omitempty" json:"meal,omitempty"`
	PolicyDRSPgind     *PolicyDRSPgind     `xml:" pgind,omitempty" json:"pgind,omitempty"`
	PolicyDRSRPH       *PolicyDRSRPH       `xml:" RPH,omitempty" json:"RPH,omitempty"`
	PolicyDRSRoutno    *PolicyDRSRoutno    `xml:" routno,omitempty" json:"routno,omitempty"`
	PolicyDRSStop      *PolicyDRSStop      `xml:" stop,omitempty" json:"stop,omitempty"`
	PolicyDRSSubid     *PolicyDRSSubid     `xml:" subid,omitempty" json:"subid,omitempty"`
	PolicyDRSTerm      *PolicyDRSTerm      `xml:" term,omitempty" json:"term,omitempty"`
	PolicyDRSTpm       *PolicyDRSTpm       `xml:" tpm,omitempty" json:"tpm,omitempty"`
	PolicyDRSWeek      *PolicyDRSWeek      `xml:" week,omitempty" json:"week,omitempty"`
}

type PolicyDRSHeaderOut struct {
	PolicyDRSSessionId *PolicyDRSSessionId `xml:" sessionId,omitempty" json:"sessionId,omitempty"`
}

type PolicyDRSNFare struct {
	PolicyDRSAgNo        *PolicyDRSAgNo        `xml:" agNo,omitempty" json:"agNo,omitempty"`
	PolicyDRSAmt         *PolicyDRSAmt         `xml:" amt,omitempty" json:"amt,omitempty"`
	PolicyDRSBkClass     *PolicyDRSBkClass     `xml:" bkClass,omitempty" json:"bkClass,omitempty"`
	PolicyDRSCarr        *PolicyDRSCarr        `xml:" carr,omitempty" json:"carr,omitempty"`
	PolicyDRSCbClass     *PolicyDRSCbClass     `xml:" cbClass,omitempty" json:"cbClass,omitempty"`
	PolicyDRSCurrCode    *PolicyDRSCurrCode    `xml:" currCode,omitempty" json:"currCode,omitempty"`
	PolicyDRSDes         *PolicyDRSDes         `xml:" des,omitempty" json:"des,omitempty"`
	PolicyDRSEI          *PolicyDRSEI          `xml:" EI,omitempty" json:"EI,omitempty"`
	PolicyDRSFbc         *PolicyDRSFbc         `xml:" fbc,omitempty" json:"fbc,omitempty"`
	PolicyDRSIb          *PolicyDRSIb          `xml:" ib,omitempty" json:"ib,omitempty"`
	PolicyDRSJourType    *PolicyDRSJourType    `xml:" jourType,omitempty" json:"jourType,omitempty"`
	PolicyDRSMaxStay     *PolicyDRSMaxStay     `xml:" maxStay,omitempty" json:"maxStay,omitempty"`
	PolicyDRSMaxStayUnit *PolicyDRSMaxStayUnit `xml:" maxStayUnit,omitempty" json:"maxStayUnit,omitempty"`
	PolicyDRSMinStay     *PolicyDRSMinStay     `xml:" minStay,omitempty" json:"minStay,omitempty"`
	PolicyDRSMinStayUnit *PolicyDRSMinStayUnit `xml:" minStayUnit,omitempty" json:"minStayUnit,omitempty"`
	PolicyDRSOb          *PolicyDRSOb          `xml:" ob,omitempty" json:"ob,omitempty"`
	PolicyDRSOri         *PolicyDRSOri         `xml:" ori,omitempty" json:"ori,omitempty"`
	PolicyDRSRPH         *PolicyDRSRPH         `xml:" RPH,omitempty" json:"RPH,omitempty"`
	PolicyDRSRuleRPH     *PolicyDRSRuleRPH     `xml:" ruleRPH,omitempty" json:"ruleRPH,omitempty"`
}

type PolicyDRSNFares struct {
	PolicyDRSNFare *PolicyDRSNFare `xml:" NFare,omitempty" json:"NFare,omitempty"`
}

type PolicyDRSOI struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSOT struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSOutput struct {
	PolicyDRSHeaderOut *PolicyDRSHeaderOut `xml:" HeaderOut,omitempty" json:"HeaderOut,omitempty"`
	PolicyDRSResult    *PolicyDRSResult    `xml:" Result,omitempty" json:"Result,omitempty"`
}

type PolicyDRSPFare struct {
	PolicyDRSAmt      *PolicyDRSAmt      `xml:" amt,omitempty" json:"amt,omitempty"`
	PolicyDRSBkClass  *PolicyDRSBkClass  `xml:" bkClass,omitempty" json:"bkClass,omitempty"`
	PolicyDRSCarr     *PolicyDRSCarr     `xml:" carr,omitempty" json:"carr,omitempty"`
	PolicyDRSCbClass  *PolicyDRSCbClass  `xml:" cbClass,omitempty" json:"cbClass,omitempty"`
	PolicyDRSCurrCode *PolicyDRSCurrCode `xml:" currCode,omitempty" json:"currCode,omitempty"`
	PolicyDRSDes      *PolicyDRSDes      `xml:" des,omitempty" json:"des,omitempty"`
	PolicyDRSFbc      *PolicyDRSFbc      `xml:" fbc,omitempty" json:"fbc,omitempty"`
	PolicyDRSOri      *PolicyDRSOri      `xml:" ori,omitempty" json:"ori,omitempty"`
	PolicyDRSRPH      *PolicyDRSRPH      `xml:" RPH,omitempty" json:"RPH,omitempty"`
	PolicyDRSRuleNo   *PolicyDRSRuleNo   `xml:" ruleNo,omitempty" json:"ruleNo,omitempty"`
	PolicyDRSRuleRPH  *PolicyDRSRuleRPH  `xml:" ruleRPH,omitempty" json:"ruleRPH,omitempty"`
}

type PolicyDRSPFares struct {
	PolicyDRSPFare []*PolicyDRSPFare `xml:" PFare,omitempty" json:"PFare,omitempty"`
}

type PolicyDRSPS struct {
	PolicyDRSDisAmt               *PolicyDRSDisAmt               `xml:" disAmt,omitempty" json:"disAmt,omitempty"`
	PolicyDRSDisCurrCode          *PolicyDRSDisCurrCode          `xml:" disCurrCode,omitempty" json:"disCurrCode,omitempty"`
	PolicyDRSEI                   *PolicyDRSEI                   `xml:" EI,omitempty" json:"EI,omitempty"`
	PolicyDRSFCs                  *PolicyDRSFCs                  `xml:" FCs,omitempty" json:"FCs,omitempty"`
	PolicyDRSFbc                  *PolicyDRSFbc                  `xml:" fbc,omitempty" json:"fbc,omitempty"`
	PolicyDRSFc                   *PolicyDRSFc                   `xml:" Fc,omitempty" json:"Fc,omitempty"`
	PolicyDRSItiType              *PolicyDRSItiType              `xml:" itiType,omitempty" json:"itiType,omitempty"`
	PolicyDRSOdType               *PolicyDRSOdType               `xml:" odType,omitempty" json:"odType,omitempty"`
	PolicyDRSOffices              *PolicyDRSOffices              `xml:" offices,omitempty" json:"offices,omitempty"`
	PolicyDRSPassengerType        *PolicyDRSPassengerType        `xml:" passengerType,omitempty" json:"passengerType,omitempty"`
	PolicyDRSRMK                  *PolicyDRSRMK                  `xml:" RMK,omitempty" json:"RMK,omitempty"`
	PolicyDRSRefundRuleIndicator  *PolicyDRSRefundRuleIndicator  `xml:" refundRuleIndicator,omitempty" json:"refundRuleIndicator,omitempty"`
	PolicyDRSReissueRuleIndicator *PolicyDRSReissueRuleIndicator `xml:" reissueRuleIndicator,omitempty" json:"reissueRuleIndicator,omitempty"`
	PolicyDRSRmkcms               *PolicyDRSRmkcms               `xml:" rmkcms,omitempty" json:"rmkcms,omitempty"`
	PolicyDRSRouts                []*PolicyDRSRouts              `xml:" Routs,omitempty" json:"Routs,omitempty"`
	PolicyDRSSeq                  *PolicyDRSSeq                  `xml:" seq,omitempty" json:"seq,omitempty"`
	PolicyDRSTaxes                *PolicyDRSTaxes                `xml:" Taxes,omitempty" json:"Taxes,omitempty"`
	PolicyDRSZValue               *PolicyDRSZValue               `xml:" zValue,omitempty" json:"zValue,omitempty"`
	PolicyDRSZValueKey            *PolicyDRSZValueKey            `xml:" zValueKey,omitempty" json:"zValueKey,omitempty"`
}

type PolicyDRSPSn struct {
	PolicyDRSPS []*PolicyDRSPS `xml:" PS,omitempty" json:"PS,omitempty"`
}

type PolicyDRSPolicy struct {
	PolicyDRSContent *PolicyDRSContent `xml:" Content,omitempty" json:"Content,omitempty"`
	PolicyDRSSeq     *PolicyDRSSeq     `xml:" Seq,omitempty" json:"Seq,omitempty"`
}

type PolicyDRSPolicyBindings struct {
	PolicyDRSPolicys *PolicyDRSPolicys `xml:" Policys,omitempty" json:"Policys,omitempty"`
}

type PolicyDRSPolicys struct {
	PolicyDRSPolicy []*PolicyDRSPolicy `xml:" Policy,omitempty" json:"Policy,omitempty"`
}

type PolicyDRSPsAvBind struct {
	PolicyDRSAvRPH   *PolicyDRSAvRPH   `xml:" avRPH,omitempty" json:"avRPH,omitempty"`
	PolicyDRSBkClass *PolicyDRSBkClass `xml:" bkClass,omitempty" json:"bkClass,omitempty"`
	PolicyDRSSeq     *PolicyDRSSeq     `xml:" seq,omitempty" json:"seq,omitempty"`
}

type PolicyDRSPsAvBinds struct {
	PolicyDRSPsAvBind []*PolicyDRSPsAvBind `xml:" PsAvBind,omitempty" json:"PsAvBind,omitempty"`
}

type PolicyDRSRMK struct {
	PolicyDRSOT *PolicyDRSOT `xml:" OT,omitempty" json:"OT,omitempty"`
}

type PolicyDRSRPH struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSRefundRules struct {
}

type PolicyDRSReissueRules struct {
}

type PolicyDRSResult struct {
	PolicyDRSError            *PolicyDRSError            `xml:" Error,omitempty" json:"Error,omitempty"`
	PolicyDRSFlightShopResult *PolicyDRSFlightShopResult `xml:" FlightShopResult,omitempty" json:"FlightShopResult,omitempty"`
}

type PolicyDRSRout struct {
	PolicyDRSArrivalAirport   *PolicyDRSArrivalAirport   `xml:" arrivalAirport,omitempty" json:"arrivalAirport,omitempty"`
	PolicyDRSArrivalDate      *PolicyDRSArrivalDate      `xml:" arrivalDate,omitempty" json:"arrivalDate,omitempty"`
	PolicyDRSArrivalTime      *PolicyDRSArrivalTime      `xml:" arrivalTime,omitempty" json:"arrivalTime,omitempty"`
	PolicyDRSBkClass          *PolicyDRSBkClass          `xml:" bkClass,omitempty" json:"bkClass,omitempty"`
	PolicyDRSCarr             *PolicyDRSCarr             `xml:" carr,omitempty" json:"carr,omitempty"`
	PolicyDRSDepartureAirport *PolicyDRSDepartureAirport `xml:" departureAirport,omitempty" json:"departureAirport,omitempty"`
	PolicyDRSDepartureDate    *PolicyDRSDepartureDate    `xml:" departureDate,omitempty" json:"departureDate,omitempty"`
	PolicyDRSDepartureTime    *PolicyDRSDepartureTime    `xml:" departureTime,omitempty" json:"departureTime,omitempty"`
	PolicyDRSFltNo            *PolicyDRSFltNo            `xml:" fltNo,omitempty" json:"fltNo,omitempty"`
	PolicyDRSOI               *PolicyDRSOI               `xml:" OI,omitempty" json:"OI,omitempty"`
}

type PolicyDRSRouts struct {
	PolicyDRSRout *PolicyDRSRout `xml:" Rout,omitempty" json:"Rout,omitempty"`
}

type PolicyDRSRule struct {
	PolicyDRSApp              *PolicyDRSApp              `xml:" app,omitempty" json:"app,omitempty"`
	PolicyDRSCancel           *PolicyDRSCancel           `xml:" cancel,omitempty" json:"cancel,omitempty"`
	PolicyDRSCombineRule      *PolicyDRSCombineRule      `xml:" CombineRule,omitempty" json:"CombineRule,omitempty"`
	PolicyDRSElig             *PolicyDRSElig             `xml:" elig,omitempty" json:"elig,omitempty"`
	PolicyDRSMaxAdv           *PolicyDRSMaxAdv           `xml:" maxAdv,omitempty" json:"maxAdv,omitempty"`
	PolicyDRSMaxAdvUnit       *PolicyDRSMaxAdvUnit       `xml:" maxAdvUnit,omitempty" json:"maxAdvUnit,omitempty"`
	PolicyDRSMinAdv           *PolicyDRSMinAdv           `xml:" minAdv,omitempty" json:"minAdv,omitempty"`
	PolicyDRSMinAdvUnit       *PolicyDRSMinAdvUnit       `xml:" minAdvUnit,omitempty" json:"minAdvUnit,omitempty"`
	PolicyDRSOth              *PolicyDRSOth              `xml:" oth,omitempty" json:"oth,omitempty"`
	PolicyDRSRPH              *PolicyDRSRPH              `xml:" RPH,omitempty" json:"RPH,omitempty"`
	PolicyDRSRebook           *PolicyDRSRebook           `xml:" rebook,omitempty" json:"rebook,omitempty"`
	PolicyDRSRuleID           *PolicyDRSRuleID           `xml:" ruleID,omitempty" json:"ruleID,omitempty"`
	PolicyDRSRuleIdNo         *PolicyDRSRuleIdNo         `xml:" ruleIdNo,omitempty" json:"ruleIdNo,omitempty"`
	PolicyDRSSeparateSaleType *PolicyDRSSeparateSaleType `xml:" separateSaleType,omitempty" json:"separateSaleType,omitempty"`
	PolicyDRSTkt              *PolicyDRSTkt              `xml:" tkt,omitempty" json:"tkt,omitempty"`
}

type PolicyDRSRules struct {
	PolicyDRSRule []*PolicyDRSRule `xml:" Rule,omitempty" json:"Rule,omitempty"`
}

type PolicyDRSSecInfo struct {
	PolicyDRSA     *PolicyDRSA     `xml:" A,omitempty" json:"A,omitempty"`
	PolicyDRSB     *PolicyDRSB     `xml:" B,omitempty" json:"B,omitempty"`
	PolicyDRSSecNo *PolicyDRSSecNo `xml:" secNo,omitempty" json:"secNo,omitempty"`
}

type PolicyDRSSeq struct {
	Text string `xml:",chardata" json:",omitempty"`
}

// type PolicyDRSSeq struct {
// 	Text string `xml:",chardata" json:",omitempty"`
// }

type PolicyDRSShopWarning struct {
	PolicyDRSCode    *PolicyDRSCode    `xml:" code,omitempty" json:"code,omitempty"`
	PolicyDRSMessage *PolicyDRSMessage `xml:" message,omitempty" json:"message,omitempty"`
}

type PolicyDRSTax struct {
	PolicyDRSAmt          *PolicyDRSAmt          `xml:" amt,omitempty" json:"amt,omitempty"`
	PolicyDRSCode         *PolicyDRSCode         `xml:" code,omitempty" json:"code,omitempty"`
	PolicyDRSCurrCode     *PolicyDRSCurrCode     `xml:" currCode,omitempty" json:"currCode,omitempty"`
	PolicyDRSTaxComponent *PolicyDRSTaxComponent `xml:" taxComponent,omitempty" json:"taxComponent,omitempty"`
}

type PolicyDRSTaxes struct {
	PolicyDRSTax []*PolicyDRSTax `xml:" Tax,omitempty" json:"Tax,omitempty"`
}

type PolicyDRSYFares struct {
	PolicyDRSYFareAmount *PolicyDRSYFareAmount `xml:" yFareAmount,omitempty" json:"yFareAmount,omitempty"`
}

type PolicyDRSAgNo struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSAirline struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSAmt struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSApp struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSArr struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSArrivalAirport struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSArrivalDate struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSArrivalTime struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSArrtm struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSAv struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSAvRPH struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSBkClass struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSCancel struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSCarr struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSCbClass struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSClass struct {
	PolicyDRSAv   *PolicyDRSAv   `xml:" av,omitempty" json:"av,omitempty"`
	PolicyDRSName *PolicyDRSName `xml:" name,omitempty" json:"name,omitempty"`
}

type PolicyDRSCode struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSCodeshare struct {
	PolicyDRSAirline *PolicyDRSAirline `xml:" airline,omitempty" json:"airline,omitempty"`
	PolicyDRSFltno   *PolicyDRSFltno   `xml:" fltno,omitempty" json:"fltno,omitempty"`
}

type PolicyDRSCombineOpenJawFlag struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSCombineRoundTripFlag struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSCurrCode struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSDep struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSDepartureAirport struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSDepartureDate struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSDepartureTime struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSDeptm struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSDes struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSDev struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSDisAmt struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSDisCurrCode struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSDt struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSElig struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSEndorsement struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSEt struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSFareBasis struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSFareCombinationFlag struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSFareRPH struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSFbc struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSFltNo struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSFltno struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSIb struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSItiType struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSJourType struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSLnk struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSMaxAdv struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSMaxAdvUnit struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSMaxStay struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSMaxStayUnit struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSMeal struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSMessage struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSMinAdv struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSMinAdvUnit struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSMinStay struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSMinStayUnit struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSName struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSOb struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSOdType struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSOffice struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSOffices struct {
	PolicyDRSOffice *PolicyDRSOffice `xml:" office,omitempty" json:"office,omitempty"`
}

type PolicyDRSOri struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSOriAmt struct {
}

type PolicyDRSOth struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSPassengerType struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSPgind struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSRebook struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSRefundRuleIndicator struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSReissueRuleIndicator struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSRmkcms struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSRoutno struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSRuleID struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSRuleIdNo struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSRuleNo struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSRuleRPH struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSSecNo struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSSeparateSaleType struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSSessionId struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSStop struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSSubid struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSSysType struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSTaxComponent struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSTerm struct {
	PolicyDRSArr *PolicyDRSArr `xml:" arr,omitempty" json:"arr,omitempty"`
	PolicyDRSDep *PolicyDRSDep `xml:" dep,omitempty" json:"dep,omitempty"`
}

type PolicyDRSTkt struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSTpm struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSWeek struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSYFareAmount struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSZValue struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type PolicyDRSZValueKey struct {
	Text string `xml:",chardata" json:",omitempty"`
}
