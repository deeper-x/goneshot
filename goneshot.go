package goneshot

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var mySigninKey = []byte(os.Getenv("SECRET_KEY"))

// GetJWT generates a JWT token
func GetJWT(client string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = client
	claims["aud"] = aud
	claims["iss"] = iss
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()

	tokenString, err := token.SignedString(mySigninKey)

	if err != nil {
		log.Println(err)
		return "", err
	}

	return tokenString, nil
}

// validateToken is the middleware function that validates the JWT token
var validateToken = func(w http.ResponseWriter, r *http.Request) (bool, error) {
	token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(("invalid Signing Method"))
		}

		// check audience claim
		checkAudience := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
		if !checkAudience {
			return nil, fmt.Errorf(("invalid aud"))
		}

		// verify iss claim
		checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
		if !checkIss {
			return nil, fmt.Errorf(("invalid iss"))
		}

		return mySigninKey, nil
	})

	if err != nil {
		return false, err
	}

	return token.Valid, nil
}

// IsAuthorized is the middleware function to check if the request is authorized
var IsAuthorized = func(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			validToken, err := validateToken(w, r)

			if err != nil {
				fmt.Fprint(w, err.Error())
			}

			if validToken {
				endpoint(w, r)
			}
		} else {
			fmt.Fprintf(w, "No Authorization Token provided")
		}
	})
}
