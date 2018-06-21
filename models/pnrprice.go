package models

type PnrPrice struct {
	PNRID           string
	Type            string
	UpPrice         float64
	UpFax           float64
	DownPrice       float64
	DownFax         float64
	UpProxyRate     float64
	UpRebatePoint   float64
	DownProxyRate   float64
	DownRebatePoint float64
	UpPoundage      float64
	DownPoundage    float64
	AllSegments     string
	FromCurrency    string
	ToCurrency      string
	Rate            float64
	YQ              float64
}
