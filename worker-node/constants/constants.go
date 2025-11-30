package constants

import (
	"fmt"
	"os"
)

var SecretKey string
var Deployment string
var AdminLoginId string
var AdminLoginPass string

func init() {
	SecretKey = os.Getenv("JWT_SECRET")
	Deployment = os.Getenv("DEPLOYMENT")

	if Deployment == "" {
		Deployment = "TEST"
	}

	if Deployment == "PRODUCTION" && SecretKey == "" {
		panic("JWT_SECRET not set in .env file in PRODUCTION")
	} else if SecretKey == "" {
		fmt.Printf("JWT_SECRET not set, the tokens would not be signed and secure")
	}

	AdminLoginId = os.Getenv("ADMIN_LOGIN")
	AdminLoginPass = os.Getenv("ADMIN_PASSWORD")
	if Deployment == "PRODUCTION" && (AdminLoginId == "" || AdminLoginPass == "") {
		panic("Admin login and password not set in production")
	}

	if AdminLoginId == "" {
		AdminLoginId = "admin"
	}
	if AdminLoginPass == "" {
		AdminLoginPass = "adminAdmin"
	}
}
