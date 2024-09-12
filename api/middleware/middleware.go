package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/lestrrat/go-jwx/jwk"
)

// JWK represents the JSON Web Key structure
type JWK struct {
	Alg string `json:"alg"`
	E   string `json:"e"`
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	N   string `json:"n"`
	Use string `json:"use"`
}

var cognitoPoolID = "COGNITO_POOL_ID"
var region = "AWS_REGION"
var jwksURL = fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", region, cognitoPoolID)

// GetCognitoPublicKeys retrieves the public keys from Cognito using the region and user pool ID
func GetCognitoPublicKeys() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		set, err := jwk.FetchHTTP(jwksURL)
		if err != nil {
			return nil, err
		}

		keyID, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("expecting JWT header to have string kid")
		}

		if key := set.LookupKeyID(keyID); len(key) == 1 {
			return key[0].Materialize()
		}

		return nil, errors.New("unable to find key")
	}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Authorization token must be Bearer token", http.StatusUnauthorized)
			return
		}

		keyFunc := GetCognitoPublicKeys()
		token, err := jwt.Parse(tokenString, keyFunc)

		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid access token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Failed to get token claims", http.StatusBadRequest)
			return
		}

		userID, ok := claims["cognito:username"].(string)
		if !ok {
			http.Error(w, "Failed to get subject claim", http.StatusBadRequest)
			return
		}

		expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
		if time.Now().After(expirationTime) {
			http.Error(w, "Access token has expired", http.StatusUnauthorized)
			return
		}

		r.Header.Add("user-id", userID)

		next.ServeHTTP(w, r)
	})
}
