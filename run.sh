#!/bin/bash

docker-compose down
docker-compose up -d
export DATABASE_URL="postgresql://postgres:postgres@localhost:5432/ktaxes?sslmode=disable"
export PORT=8080
export ADMIN_USERNAME=adminTax
export ADMIN_PASSWORD=admin!
go run main.go
