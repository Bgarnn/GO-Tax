package service

import (
	"math"
	"net/http"

	"github.com/Bgarnn/assessment-tax/database"
	handler "github.com/Bgarnn/assessment-tax/struct"
	"github.com/labstack/echo"
)

type Err struct {
	Message string `json:"message"`
}

type Level struct {
	Level             int
	LevelString       string
	MinAmount         float64
	MaxAmount         float64
	TaxRatePercentage float64
}

func CreateLevels() []Level {
	return []Level{
		{Level: 1, LevelString: "0 - 150,000", MinAmount: 0.0, MaxAmount: 150000, TaxRatePercentage: 0.00},
		{Level: 2, LevelString: "150,001 - 500,000", MinAmount: 150001.00, MaxAmount: 500000.00, TaxRatePercentage: 0.10},
		{Level: 3, LevelString: "500,001 - 1,000,000", MinAmount: 500001.00, MaxAmount: 1000000.00, TaxRatePercentage: 0.15},
		{Level: 4, LevelString: "1,000,001 - 2,000,000", MinAmount: 1000001.00, MaxAmount: 2000000.00, TaxRatePercentage: 0.20},
		{Level: 5, LevelString: "2,000,001 ขึ้นไป", MinAmount: 2000001.00, MaxAmount: math.MaxFloat64, TaxRatePercentage: 0.35},
	}
}

func Calculate(c echo.Context, data database.DataStruct) error {
	var request handler.RequestCalculation
	if err := c.Bind(&request); err != nil {
		return err
	}
	request.Wht = ValidateWht(request.Wht, request.TotalIncome)
	if request.Wht == -1 {
		return c.JSON(http.StatusBadRequest, Err{Message: "Invalid wht"})
	}
	taxableIncome, err := AllowanceCalculate(data, request)
	if err != nil {
		return err
	}
	taxAmount, taxLevels := TaxLevelCalculate(taxableIncome)
	var taxRefund float64
	taxRefund, taxAmount = WhtCalculate(request.Wht, taxAmount)
	response := handler.ResponseCalculation{TaxRefund: taxRefund, Tax: taxAmount, TaxLevel: taxLevels}
	return c.JSON(http.StatusOK, response)
}

func WhtCalculate(wht, taxAmount float64) (float64, float64) {
	var taxRefund float64
	taxAmount -= wht
	if taxAmount < 0 {
		taxRefund = (-1) * taxAmount
		taxAmount = 0
		return taxRefund, taxAmount
	}
	return 0, taxAmount
}

func AllowanceCalculate(data database.DataStruct, request handler.RequestCalculation) (float64, error) {
	var totalAllowanceAmount float64
	DonationValidate(&request)
	for _, a := range request.Allowances {
		totalAllowanceAmount += a.Amount
	}
	totalAllowanceAmount += data.PersonalAllowance
	taxableIncome := request.TotalIncome - totalAllowanceAmount
	return taxableIncome, nil
}

func TaxLevelCalculate(taxableIncome float64) (float64, []handler.TaxLevelArr) {
	var taxLevelsArr []handler.TaxLevelArr
	var taxResultTotal float64
	taxLevelDetail := CreateLevels()

	levelOfTax := GetTaxLevel(taxableIncome, taxLevelDetail)
	for i := levelOfTax; i >= 0; i-- {
		totalIncomeThisLevel := min(taxableIncome, taxLevelDetail[i].MaxAmount) - taxLevelDetail[i].MinAmount + 1
		taxResultThisLevel := totalIncomeThisLevel * taxLevelDetail[i].TaxRatePercentage

		newTaxLevel := handler.TaxLevelArr{Level: taxLevelDetail[i].LevelString, Tax: taxResultThisLevel}
		taxLevelsArr = append([]handler.TaxLevelArr{newTaxLevel}, taxLevelsArr...)

		taxResultTotal += taxResultThisLevel
		taxableIncome -= totalIncomeThisLevel
	}
	for i := levelOfTax + 1; i <= 4; i++ {
		taxLevelsArr = append(taxLevelsArr, handler.TaxLevelArr{Level: taxLevelDetail[i].LevelString, Tax: 0})
	}
	return taxResultTotal, taxLevelsArr
}

func GetTaxLevel(taxableIncome float64, levels []Level) int {
	for _, level := range levels {
		if taxableIncome >= level.MinAmount && taxableIncome <= level.MaxAmount {
			return level.Level - 1
		}
	}
	return -1
}

func ValidateWht(amount, totalIncome float64) float64 {
	if amount > totalIncome || amount < 0 {
		return -1
	}
	return amount
}

func DonationValidate(request *handler.RequestCalculation) {
	for i := range request.Allowances {
		if request.Allowances[i].AllowanceType == "donation" && request.Allowances[i].Amount > 100000 {
			request.Allowances[i].Amount = 100000
		}
	}
}
