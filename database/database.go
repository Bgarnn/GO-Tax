package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"os"

	handler "github.com/Bgarnn/assessment-tax/struct"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
)

type Err struct {
	Message string `json:"message"`
}

type DataStruct struct {
	PersonalAllowance float64
	MaxKReceipt       float64
}

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database failed", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("Ping database failed", err)
	}
	createTb := `CREATE TABLE IF NOT EXISTS allowance ( id SERIAL PRIMARY KEY, personal FLOAT, maxKReceipt FLOAT);`
	_, err = DB.Exec(createTb)
	if err != nil {
		log.Fatal("Create table failed", err)
	}
	DB.QueryRow("INSERT INTO allowance (personal, maxKReceipt) values ($1, $2) RETURNING id", 60000, 50000)
}

func GetPersonal(db *sql.DB) (float64, error) {
	var personalAllowance float64
	err := db.QueryRow("SELECT personal FROM allowance WHERE id = $1", 1).Scan(&personalAllowance)
	if err != nil {
		return 0, fmt.Errorf("GetPersonal failed: %v", err)
	}
	personalAllowance = ValidatePersonal(personalAllowance)
	return personalAllowance, nil
}

func ValidatePersonal(amount float64) float64 {
	if amount > 100000 {
		return (100000)
	} else if amount <= 10000 {
		return (10001)
	} else {
		return (amount)
	}
}

func UpdatePersonal(c echo.Context, data DataStruct) error {
	db := DB
	var request handler.RequestDeduction
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	personalAllowance := ValidatePersonal(request.Amount)
	data.PersonalAllowance = request.Amount
	stmt, err := db.Prepare(`UPDATE allowance SET personal = $1`)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	if _, err := stmt.Exec(personalAllowance); err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	response := map[string]float64{"personalDeduction": personalAllowance}
	return c.JSON(http.StatusOK, response)
}

func GetMaxKReceipt(db *sql.DB) (float64, error) {
	var maxKReceiptAllowance float64
	err := db.QueryRow("SELECT maxKReceipt FROM allowance WHERE id = $1", 1).Scan(&maxKReceiptAllowance)
	if err != nil {
		return 0, fmt.Errorf("GetMaxKReceipt failed: %v", err)
	}
	maxKReceiptAllowance = ValidateMaxKReceipt(maxKReceiptAllowance)
	return maxKReceiptAllowance, nil
}

func ValidateMaxKReceipt(amount float64) float64 {
	if amount > 100000 {
		return (100000)
	} else if amount <= 0 {
		return (1)
	} else {
		return (amount)
	}
}

func UpdateMaxKReceipt(c echo.Context, data DataStruct) error {
	db := DB
	var request handler.RequestDeduction
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	maxKReceipt := ValidateMaxKReceipt(request.Amount)
	data.MaxKReceipt = request.Amount
	stmt, err := db.Prepare(`UPDATE allowance SET maxKReceipt = $1`)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	if _, err := stmt.Exec(maxKReceipt); err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	response := map[string]float64{"kReceipt": maxKReceipt}
	return c.JSON(http.StatusOK, response)
}
