package handler

type RequestCalculation struct {
	TotalIncome float64         `json:"totalIncome"`
	Wht         float64         `json:"wht"`
	Allowances  []AllowancesArr `json:"allowances"`
}

type AllowancesArr struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type ResponseCalculation struct {
	TaxRefund float64       `json:"taxRefund,omitempty"`
	Tax       float64       `json:"tax"`
	TaxLevel  []TaxLevelArr `json:"taxLevel,omitempty"`
}

type TaxLevelArr struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

type RequestDeduction struct {
	Amount float64 `json:"amount"`
}

type RequestCsv struct {
	TotalIncome float64 `json:"totalIncome"`
	Tax         float64 `json:"tax"`
}

type ResponseCSV struct {
	TotalIncome float64 `json:"totalIncome"`
	Tax         float64 `json:"tax"`
	TaxRefund   float64 `json:"taxRefund,omitempty"`
}
