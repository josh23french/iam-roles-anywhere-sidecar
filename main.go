package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	helper "github.com/aws/rolesanywhere-credential-helper/aws_signing_helper"
)

func getBoolEnv(key string) bool {
	val, err := strconv.ParseBool(os.Getenv(key))
	if err != nil {
		return false
	}
	return val
}

func getIntEnv(key string, defalt int) int {
	val, err := strconv.ParseInt(os.Getenv(key), 10, 0)
	if err != nil {
		return defalt
	}
	return int(val)

}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/creds", func(c echo.Context) error {
		credentialsOptions := helper.CredentialsOpts{
			PrivateKeyId:        os.Getenv("PRIVATE_KEY_ID"),
			CertificateId:       os.Getenv("CERTIFICATE_ID"),
			CertificateBundleId: os.Getenv("CERTIFICATE_BUNDLE_ID"),
			RoleArn:             os.Getenv("ROLE_ARN"),
			ProfileArnStr:       os.Getenv("PROFILE_ARN"),
			TrustAnchorArnStr:   os.Getenv("TRUST_ANCHOR_ID"),
			SessionDuration:     getIntEnv("SESSION_DURATION", 3600),
			Region:              os.Getenv("AWS_REGION"),
			Endpoint:            os.Getenv("ENDPOINT"),
			NoVerifySSL:         getBoolEnv("NO_VERIFY_SSL"),
			WithProxy:           getBoolEnv("WITH_PROXY"),
			Debug:               getBoolEnv("DEBUG"),
			Version:             os.Getenv("CREDENTIAL_VERSION"),
		}
		output, err := helper.GenerateCredentials(&credentialsOptions)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, &ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusOK, output)
	})

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start("[::1]:" + httpPort))
}
