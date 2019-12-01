package authorization

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/TheStevbeef/communityServiceGo/utils"
	jwt "github.com/dgrijalva/jwt-go"
)

type Keys struct {
	Keys []Key `json:"keys"`
}
type Key struct {
	Kty string `json:"kty"`
	E   string `json:"E"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	N   string `json:"n"`
}

func RequireTokenAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, error := jwt.Parse(bearerToken[1], getKey)
				if error != nil {
					utils.RespondWithError(w, http.StatusInternalServerError, error.Error())
					return
				}
				if token.Valid {
					//context.Set(r, "decoded", token.Claims)
					next.ServeHTTP(w, r)
				} else {
					utils.RespondWithError(w, http.StatusUnauthorized, "Invalid authorization token")
				}
			}
		} else {
			utils.RespondWithError(w, http.StatusUnauthorized, "An authorization header is required")
		}

	})
}

func getKey(token *jwt.Token) (interface{}, error) {
	jsonFile, err := os.Open("keys.json")
	if err != nil {
		return nil, fmt.Errorf("Could not read keys.json File")
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("Could not read keys.json File")
	}
	var keys Keys
	json.Unmarshal(byteValue, &keys)
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("There was an error")
	}
	keyID, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("expecting JWT header to have string kid")
	}
	for _, key := range keys.Keys {
		if key.Kid == keyID {
			return `    {
				"kty": "RSA",
				"e": "AQAB",
				"use": "sig",
				"kid": "1",
				"alg": "RS256",
				"n": "lSrbibw4dTDY-RKj1jNYoRSRUi6ZUnc1mqe9a9x4QHjlc2Mjxf3Hc5SwmuWo4lnhQ1_CjYwZElFyZp0ajXhyTGasT9LLblE6iYEk4xXRl8avbM5LxAlJUY4viwQsQggjNhPY9dLjQuS0RZ7VXw0egRs8IfInR5CWSG_8NF1JZLmQtQYkqKs-UjAdA6mJ2BU382iFVx6KdYQ6ondKLOEu7GC0wbaeyjk7dxGeuFRpNUSShyBVggtWjmxmjS9cnNrFOP-eBkggeOuuMSl-HDUTAWoW6sgNykezzgBvtbUUxeBo3Rpdsflt6eO-bB-mwRhFkGF_M8SgGNYwvUT9ck6mEQ"
			  }`, nil
		}
	}
	return nil, fmt.Errorf("unable to find key %q", keyID)
}
