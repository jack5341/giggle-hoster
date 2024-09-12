package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/jack5341/giggle-hoster/internal/cognito"
)

var (
	ErrCogniteCouldNotBeInitiliazed          = errors.New("auth provider could not be initiliazed")
	ErrEmailUsernamePasswordIsRequiredInputs = errors.New("email username or password is required inputs")
	ErrEmailAndPasswordIsRequired            = errors.New("email and password is required")
)

type User struct {
	Username string `form:"username"`
	Email    string `form:"email"`
	Password string `form:"password"`
	Code     string `form:"code"`
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if u.Email == "" || u.Username == "" || u.Password == "" {
		http.Error(w, ErrEmailUsernamePasswordIsRequiredInputs.Error(), http.StatusBadRequest)
		return
	}

	c, err := cognito.NewCognitoClient()
	if err != nil {
		http.Error(w, errors.Join(ErrCogniteCouldNotBeInitiliazed, err).Error(), http.StatusInternalServerError)
		return
	}

	input := cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(os.Getenv("COGNITO_CLIENT_ID")),
		Password: aws.String(u.Password),
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("username"),
				Value: &u.Username,
			},
		},
		Username: aws.String(u.Email),
	}

	_, responseCode, err := c.SignUp(&input)
	if err != nil {
		http.Error(w, err.Error(), responseCode)
		return
	}

	w.WriteHeader(responseCode)
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if u.Email == "" || u.Password == "" {
		http.Error(w, ErrEmailAndPasswordIsRequired.Error(), http.StatusBadRequest)
		return
	}

	c, err := cognito.NewCognitoClient()
	if err != nil {
		http.Error(w, errors.Join(ErrCogniteCouldNotBeInitiliazed, err).Error(), http.StatusInternalServerError)
		return
	}

	input := cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String(cognitoidentityprovider.AuthFlowTypeUserPasswordAuth),
		ClientId: aws.String(os.Getenv("COGNITO_CLIENT_ID")),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(u.Email),
			"PASSWORD": aws.String(u.Password),
		},
	}

	output, responseCode, err := c.InitateAuth(&input)
	if err != nil {
		http.Error(w, err.Error(), responseCode)
		return
	}

	refreshToken := http.Cookie{
		Name:     "refreshToken",
		Value:    *output.AuthenticationResult.RefreshToken,
		HttpOnly: true,
		Secure:   true,
	}

	token := http.Cookie{
		Name:     "token",
		Value:    *output.AuthenticationResult.IdToken,
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, &refreshToken)
	http.SetCookie(w, &token)

	w.WriteHeader(responseCode)
}
