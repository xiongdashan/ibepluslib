package models

type Gender string

const (
	Man      Gender = "MALE"
	Feminine Gender = "FEMALE"
)

/*********************************/
//PersonType 乘客类型
type PersonType string

const (
	Adult  PersonType = "ADT"
	Child  PersonType = "CHD"
	Infant PersonType = "INF"
)

/************************************/
//IdCardNo 证件类型
type IdCardNo string

const (
	IDCard           IdCardNo = "NI" //身份证
	Passport         IdCardNo = "PP" //港澳通行证
	Other            IdCardNo = "2"
	HometownCard     IdCardNo = "3"
	MilitaryLicense  IdCardNo = "4"
	PoliceOfficerPas IdCardNo = "5"
	HKAndMacPas      IdCardNo = "6"
	MTPs             IdCardNo = "7"
	TaiwanPass       IdCardNo = "8"
	PermanentCard    IdCardNo = "9"
)

/************************************/

type TravelType string

const (
	OW TravelType = "OW" //单程
	RT TravelType = "RT" //往返
)

type TripType string

const (
	OB TripType = "OB" //去程
	IB TripType = "IB" //回程
)
