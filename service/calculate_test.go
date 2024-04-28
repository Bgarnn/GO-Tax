package service

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Bgarnn/assessment-tax/database"
	handler "github.com/Bgarnn/assessment-tax/struct"
	"github.com/labstack/echo"
)

var LevelDetail = []Level{
	{Level: 1, LevelString: "0 - 150,000", MinAmount: 0.0, MaxAmount: 150000, TaxRatePercentage: 0.00},
	{Level: 2, LevelString: "150,001 - 500,000", MinAmount: 150001.00, MaxAmount: 500000.00, TaxRatePercentage: 0.10},
	{Level: 3, LevelString: "500,001 - 1,000,000", MinAmount: 500001.00, MaxAmount: 1000000.00, TaxRatePercentage: 0.15},
	{Level: 4, LevelString: "1,000,001 - 2,000,000", MinAmount: 1000001.00, MaxAmount: 2000000.00, TaxRatePercentage: 0.20},
	{Level: 5, LevelString: "2,000,001 ขึ้นไป", MinAmount: 2000001.00, MaxAmount: math.MaxFloat64, TaxRatePercentage: 0.35},
}

func TestCreateLevels(t *testing.T) {
	expectedLevels := LevelDetail

	actualLevels := CreateLevels()

	for i, expected := range expectedLevels {
		actual := actualLevels[i]
		if actual.Level != expected.Level ||
			actual.LevelString != expected.LevelString ||
			actual.MinAmount != expected.MinAmount ||
			actual.MaxAmount != expected.MaxAmount ||
			actual.TaxRatePercentage != expected.TaxRatePercentage {
			t.Errorf("Incorrect at index %d. Expected: %+v, Got: %+v", i, expected, actual)
		}
	}
}

func TestCalculate(t *testing.T) {
	t.Run("should return 29000 with status 200", func(t *testing.T) {
		body, err := json.Marshal(handler.RequestCalculation{
			TotalIncome: 500000.0,
			Wht:         0,
			Allowances: []handler.AllowancesArr{
				{AllowanceType: "donation", Amount: 0.0},
			},
		})
		if err != nil {
			t.Errorf("Create request failed: %v", err)
		}
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, res)
		data := database.DataStruct{
			PersonalAllowance: 60000.0,
			MaxKReceipt:       50000.0,
		}

		Calculate(c, data)

		if res.Result().StatusCode != http.StatusOK {
			t.Errorf("expected status %v but got status %v", http.StatusOK, res.Result().StatusCode)
		}
		want := handler.ResponseCalculation{
			Tax: 29000.0,
			TaxLevel: []handler.TaxLevelArr{
				{Level: "0 - 150,000", Tax: 0.00},
				{Level: "150,001 - 500,000", Tax: 29000.00},
				{Level: "500,001 - 1,000,000", Tax: 0.00},
				{Level: "1,000,001 - 2,000,000", Tax: 0.00},
				{Level: "2,000,001 ขึ้นไป", Tax: 0.00},
			},
		}
		var got handler.ResponseCalculation
		if err := json.Unmarshal(res.Body.Bytes(), &got); err != nil {
			t.Errorf("Cannot unmarshal json: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v but got %v", want, got)
		}
	})
}

func TestWhtCalculate(t *testing.T) {
	t.Run("should refund 5000 and tax 0", func(t *testing.T) {
		wht, taxAmount := 30000.0, 25000.0

		actualTaxRefund, actualTaxAmount := WhtCalculate(wht, taxAmount)

		if !reflect.DeepEqual(actualTaxRefund, 5000.0) {
			t.Errorf("expected %v but got %v", 5000.0, actualTaxRefund)
		}
		if !reflect.DeepEqual(actualTaxAmount, 0.0) {
			t.Errorf("expected %v but got %v", 0.0, actualTaxAmount)
		}
	})
	t.Run("should refund 0 and tax 5000", func(t *testing.T) {
		wht, taxAmount := 25000.0, 30000.0

		actualTaxRefund, actualTaxAmount := WhtCalculate(wht, taxAmount)

		if !reflect.DeepEqual(actualTaxRefund, 0.0) {
			t.Errorf("expected %v but got %v", 0.0, actualTaxRefund)
		}
		if !reflect.DeepEqual(actualTaxAmount, 5000.0) {
			t.Errorf("expected %v but got %v", 5000.0, actualTaxAmount)
		}
	})
}

func TestTaxLevelCalculate(t *testing.T) {
	t.Run("should return 1", func(t *testing.T) {
		taxableIncome := 500000.0
		wantTaxLevelsArr := []handler.TaxLevelArr{
			{Level: "0 - 150,000", Tax: 0.00},
			{Level: "150,001 - 500,000", Tax: 35000.00},
			{Level: "500,001 - 1,000,000", Tax: 0.00},
			{Level: "1,000,001 - 2,000,000", Tax: 0.00},
			{Level: "2,000,001 ขึ้นไป", Tax: 0.00},
		}

		gotTaxResultTotal, gotTaxLevelsArr := TaxLevelCalculate(taxableIncome)

		if !reflect.DeepEqual(wantTaxLevelsArr, gotTaxLevelsArr) {
			t.Errorf("expected %v but got %v", wantTaxLevelsArr, gotTaxLevelsArr)
		}
		if !reflect.DeepEqual(35000.00, gotTaxResultTotal) {
			t.Errorf("expected %v but got %v", 35000.00, gotTaxResultTotal)
		}
	})
}

func TestGetTaxLevel(t *testing.T) {
	t.Run("should return 1", func(t *testing.T) {
		level := LevelDetail

		got := GetTaxLevel(500000, level)

		if !reflect.DeepEqual(got, 1) {
			t.Errorf("expected %v but got %v", 1, got)
		}
	})
	t.Run("should return -1", func(t *testing.T) {
		level := LevelDetail

		got := GetTaxLevel(-1, level)

		if !reflect.DeepEqual(got, -1) {
			t.Errorf("expected %v but got %v", -1, got)
		}
	})

}

func TestValidateWht(t *testing.T) {
	t.Run("should return amount", func(t *testing.T) {
		amount, totalIncome := 30000.0, 35000.0
		want := amount

		got := ValidateWht(amount, totalIncome)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v but got %v", want, got)
		}
	})
	t.Run("should return -1", func(t *testing.T) {
		amount, totalIncome := 35000.0, 30000.0

		got := ValidateWht(amount, totalIncome)

		if !reflect.DeepEqual(got, -1.0) {
			t.Errorf("expected %v but got %v", -1.0, got)
		}
	})
}

func TestDonationValidate(t *testing.T) {
	t.Run("should return 100000", func(t *testing.T) {
		request := &handler.RequestCalculation{
			TotalIncome: 500000.0,
			Wht:         0.0,
			Allowances: []handler.AllowancesArr{
				{AllowanceType: "donation", Amount: 150000.0},
			},
		}

		DonationValidate(request)

		if request.Allowances[0].Amount != 100000.0 {
			t.Errorf("expected %f, but got %f", 100000.0, request.Allowances[0].Amount)
		}
	})
}

func TestKReceiptValidate(t *testing.T) {
	t.Run("should return 100000", func(t *testing.T) {
		request := &handler.RequestCalculation{
			TotalIncome: 500000.0,
			Wht:         0.0,
			Allowances: []handler.AllowancesArr{
				{AllowanceType: "k-receipt", Amount: 150000.0},
			},
		}
		data := database.DataStruct{
			PersonalAllowance: 60000.0,
			MaxKReceipt:       50000.0,
		}

		KReceiptValidate(data, request)

		if request.Allowances[0].Amount != data.MaxKReceipt {
			t.Errorf("expected %f, but got %f", data.MaxKReceipt, request.Allowances[0].Amount)
		}
	})
}
