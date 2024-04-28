package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Bgarnn/assessment-tax/database"
	"github.com/Bgarnn/assessment-tax/service"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	var err error
	data := database.DataStruct{
		PersonalAllowance: 60000.0,
		MaxKReceipt:       50000.0,
	}
	database.Init()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/tax/calculations", func(c echo.Context) error {
		data, err = UpdateData(data)
		if err != nil {
			return err
		}
		return service.Calculate(c, data)
	})
	e.POST("/tax/calculations/upload-csv", func(c echo.Context) error {
		data, err = UpdateData(data)
		if err != nil {
			return err
		}
		return service.Csv(c, data)
	})

	g := e.Group("/admin")
	g.Use(middleware.BasicAuth(AuthMiddleware))
	g.POST("/deductions/personal", func(c echo.Context) error {
		return database.UpdatePersonal(c, data)
	})

	go func() {
		if err := e.Start(":" + os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	fmt.Println("\nshutting down the server")
}

func AuthMiddleware(username, password string, c echo.Context) (bool, error) {
	expectedUsername := os.Getenv("ADMIN_USERNAME")
	expectedPassword := os.Getenv("ADMIN_PASSWORD")

	if username == expectedUsername && password == expectedPassword {
		return true, nil
	}
	return false, nil
}

func UpdateData(data database.DataStruct) (database.DataStruct, error) {
	var err error
	data.PersonalAllowance, err = database.GetPersonal(database.DB)
	if err != nil {
		return data, fmt.Errorf("GetPersonal error: %v", err)
	}
	data.MaxKReceipt, err = database.GetMaxKReceipt(database.DB)
	if err != nil {
		return data, fmt.Errorf("GetMaxKReceipt error: %v", err)
	}
	return data, nil
}
