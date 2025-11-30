package constants

import "os"

var SecretKey string
var Deployment string

func init() {
	SecretKey = os.Getenv("JWT_SECRET")
	Deployment = os.Getenv("DEPLOYMENT")
	if Deployment == "" {
		Deployment = "TEST"
	}
}
