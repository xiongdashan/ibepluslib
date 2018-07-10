package models

import (
	"fmt"
	"strings"
)

//Traveler 乘客信息.....
type Traveler struct {
	PersonName     string
	Gender         string
	Type           string
	IDCardType     string
	IDCardNo       string
	Birthday       string
	Nationality    string
	IDIssueCountry string
	IDIssueDate    string
	IDExpireDate   string
	OrderID        string
	PNRID          string
	Text           string
	TravelerID     string
	GivenName      string
	SurName        string
	Mobile         string
	TicketNo       string
	// Price          float64
	// Fax            float64
	RPH          string
	Airline      string
	CombinStatus string
}

func (t *Traveler) GetRPHP() string {
	if strings.HasPrefix(t.RPH, "P") {
		return t.RPH
	}
	return "P" + t.RPH
}

func (t *Traveler) GetDocs() string {
	ary := []string{
		t.IDCardType,
		t.Nationality,
		t.IDCardNo,
		t.IDIssueCountry,
		t.Birthday,
		t.Gender,
		t.IDExpireDate,
		t.PersonName,
		t.GetRPHP(),
	}

	//"DOCS UA HK1 P/CN/EB8660653/CN/25JUL94/M/27MAY28/YANG/HU/P7
	ssrT := strings.Join(ary, "/")
	all := fmt.Sprintf("DOCS %s %s1 %s", t.Airline, t.CombinStatus, ssrT)
	return all
}

func (t *Traveler) GetFOID() string {
	all := fmt.Sprintf("FOID %s %s1 NI%s/%s", t.Airline, t.CombinStatus, t.IDCardNo, t.GetRPHP())
	return all
}

type ITravelerSSR interface {
	Check(traveler *Traveler, personname, ssr string) bool
}

type TravlerDocs struct {
}

func (d *TravlerDocs) Check(traveler *Traveler, personname, ssr string) bool {
	docs := traveler.GetDocs()
	return ssr == docs
}

type TravlerFOID struct {
}

func (f *TravlerFOID) Check(traveler *Traveler, personname, ssr string) bool {

	foid := traveler.GetFOID()

	if foid == ssr {
		if personname != traveler.PersonName {
			fmt.Printf("姓名不匹配：%s-%s\n", traveler.PersonName, personname)
		}
	}

	// fmt.Println(ssr)
	// fmt.Println(foid)
	// fmt.Println("***************************")

	return foid == ssr && personname == traveler.PersonName

}
