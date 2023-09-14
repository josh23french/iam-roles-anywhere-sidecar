package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	helper "github.com/aws/rolesanywhere-credential-helper/aws_signing_helper"
)

// https://github.com/awslabs/amazon-ecs-local-container-endpoints/blob/ce24b29f9c7e880f2b7bfc285d816dbc0d06c499/local-container-endpoints/handlers/types.go
// for SOME reason, the signing helper doesn't give the same shape expected from the container creds url 🤯

// CredentialResponse is used to marshal the JSON response for the Credentials Service
type CredentialResponse struct {
	AccessKeyID     string `json:"AccessKeyId"`
	Expiration      string
	RoleArn         string
	SecretAccessKey string
	Token           string
}

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
	log.SetFlags(log.LUTC | log.Ldate | log.Ltime | log.Lmicroseconds)

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

	listen := os.Getenv("LISTEN")
	if listen == "" {
		// default to localhost:8080, but allow it to be overridden...
		// specifically for docker on mac where ipv6 somehow *still* isn't supported 🤯
		listen = "[::1]:8080"
	}

	http.HandleFunc("/creds", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		output, err := helper.GenerateCredentials(&credentialsOptions)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&ErrorResponse{Error: err.Error()})
			return
		}
		json.NewEncoder(w).Encode(&CredentialResponse{
			AccessKeyID:     output.AccessKeyId,
			Expiration:      output.Expiration,
			RoleArn:         credentialsOptions.RoleArn,
			SecretAccessKey: output.SecretAccessKey,
			Token:           output.SessionToken,
		})
	})

	http.HandleFunc("/healthz", healthzHandler)

	log.Printf("listening on %v\n", listen)
	log.Fatal(http.ListenAndServe(listen, requestLogger(http.DefaultServeMux)))
}

func healthzHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 OK"))
}
