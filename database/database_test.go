package database

import (
	"database/sql"
	"os"
	"testing"
)

func TestValidatePersonal(t *testing.T) {
	t.Run("should return 100000", func(t *testing.T) {
		amount := 200000

		got := ValidatePersonal(float64(amount))

		if got != 100000.0 {
			t.Errorf("expected %f, but got %f", 100000.0, got)
		}
	})
	t.Run("should return 15000", func(t *testing.T) {
		amount := 15000

		got := ValidatePersonal(float64(amount))

		if got != 15000.0 {
			t.Errorf("expected %f, but got %f", 15000.0, got)
		}
	})
	t.Run("should return 10001", func(t *testing.T) {
		amount := 200

		got := ValidatePersonal(float64(amount))

		if got != 10001.0 {
			t.Errorf("expected %f, but got %f", 10001.0, got)
		}
	})
}

func TestValidateMaxKReceipt(t *testing.T) {
	t.Run("should return 100000", func(t *testing.T) {
		amount := 200000

		got := ValidateMaxKReceipt(float64(amount))

		if got != 100000.0 {
			t.Errorf("expected %f, but got %f", 100000.0, got)
		}
	})
	t.Run("should return 1", func(t *testing.T) {
		amount := -15000

		got := ValidateMaxKReceipt(float64(amount))

		if got != 1.0 {
			t.Errorf("expected %f, but got %f", 1.0, got)
		}
	})
	t.Run("should return 1500", func(t *testing.T) {
		amount := 1500

		got := ValidateMaxKReceipt(float64(amount))

		if got != 1500.0 {
			t.Errorf("expected %f, but got %f", 1500.0, got)
		}
	})
}

func TestInit(t *testing.T) {
	os.Setenv("DATABASE_URL", "postgresql://postgres:postgres@localhost:5432/ktaxes?sslmode=disable")

	Init()

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()

	var tableExists bool
	err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'allowance')").Scan(&tableExists)
	if err != nil {
		t.Fatalf("Failed to check table existence: %v", err)
	}
	if !tableExists {
		t.Error("Table 'allowance' does not exist")
	}

	var columns []string
	rows, err := db.Query("SELECT column_name FROM information_schema.columns WHERE table_name = 'allowance'")
	if err != nil {
		t.Fatalf("Failed to query table columns: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var column string
		if err := rows.Scan(&column); err != nil {
			t.Fatalf("Failed to scan column name: %v", err)
		}
		columns = append(columns, column)
	}
	expectedColumns := []string{"id", "personal", "maxkreceipt"}
	for _, expectedColumn := range expectedColumns {
		found := false
		for _, column := range columns {
			if column == expectedColumn {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected column %s not found in table 'allowance'", expectedColumn)
		}
	}
}
