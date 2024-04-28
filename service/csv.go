package service

import (
	"encoding/csv"
	"log"
	"net/http"
	"strconv"

	"github.com/Bgarnn/assessment-tax/database"
	handler "github.com/Bgarnn/assessment-tax/struct"
	"github.com/labstack/echo"
)

type Err struct {
	Message string `json:"message"`
}

func Csv(c echo.Context, dt database.DataStruct) error {
	var response []handler.ResponseCSV
	file, err := c.FormFile("taxFile")
	if err != nil {
		return c.JSON(http.StatusNotFound, Err{Message: "taxFile(FormFile) error"})
	}
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "taxFile(fileOpen) error"})
	}
	defer src.Close()
	reader := csv.NewReader(src)
	data, err := reader.ReadAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "taxFile(ReadAll) error"})
	}
	for i, record := range data {
		if i == 0 {
			continue
		}
		totalIncome, wht, donation, err := ParseData(record)
		if err != nil {
			return c.JSON(http.StatusBadRequest, Err{Message: "ParseData error"})
		}
		taxableIncome := totalIncome - dt.PersonalAllowance - donation
		taxAmount, _ := TaxLevelCalculate(taxableIncome)
		var taxRefund float64
		taxRefund, taxAmount = WhtCalculate(wht, taxAmount)
		response = append(response, handler.ResponseCSV{TotalIncome: totalIncome, Tax: taxAmount, TaxRefund: taxRefund})
	}
	return c.JSON(http.StatusOK, response)
}

func ParseData(record []string) (float64, float64, float64, error) {
	totalIncome, err := strconv.ParseFloat(record[0], 64)
	if err != nil {
		log.Printf("Invalid totalIncome: %v, error: %v", record[0], err)
		return 0, 0, 0, err
	}
	wht, err := strconv.ParseFloat(record[1], 64)
	if err != nil {
		log.Printf("Invalid wht: %v, error: %v", record[1], err)
		return 0, 0, 0, err
	}
	wht = ValidateWht(wht, totalIncome)
	if wht == -1 {
		log.Printf("Invalid wht")
		return 0, 0, 0, err
	}
	donation, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		log.Printf("Invalid donation: %v, error: %v", record[2], err)
		return 0, 0, 0, err
	}
	donation = ValidateDonation(donation)
	return totalIncome, wht, donation, nil
}

func ValidateDonation(amount float64) float64 {
	if amount > 100000 {
		return 100000
	}
	return amount
}
