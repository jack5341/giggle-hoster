package cognito

import (
	"errors"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type Cognito struct {
	Client *cognitoidentityprovider.CognitoIdentityProvider
}

func NewCognitoClient() (Cognito, error) {
	cog := new(Cognito)

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_COGNITO_REGION")),
	}))

	cognitoClient := cognitoidentityprovider.New(sess)
	cog.Client = cognitoClient

	return *cog, nil
}

func (c *Cognito) SignUp(input *cognitoidentityprovider.SignUpInput) (*cognitoidentityprovider.SignUpOutput, int, error) {
	output, err := c.Client.SignUp(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cognitoidentityprovider.ErrCodeInvalidParameterException:
				return &cognitoidentityprovider.SignUpOutput{}, http.StatusBadRequest, errors.New("an account with the given email already exists")
			case cognitoidentityprovider.ErrCodeUsernameExistsException:
				return &cognitoidentityprovider.SignUpOutput{}, http.StatusConflict, errors.New("username is already exist")
			case cognitoidentityprovider.ErrCodeInvalidPasswordException:
				return &cognitoidentityprovider.SignUpOutput{}, http.StatusBadRequest, errors.New("password must include uppercase, special-character and number")
			default:
				return &cognitoidentityprovider.SignUpOutput{}, http.StatusInternalServerError, errors.New("something went wrong while sign up")
			}
		}
	}

	return output, http.StatusOK, nil
}

func (c *Cognito) InitateAuth(input *cognitoidentityprovider.InitiateAuthInput) (*cognitoidentityprovider.InitiateAuthOutput, int, error) {
	authOutput, err := c.Client.InitiateAuth(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cognitoidentityprovider.ErrCodeUserNotConfirmedException:
				return &cognitoidentityprovider.InitiateAuthOutput{}, http.StatusBadRequest, errors.New("email is not confirmed")
			case cognitoidentityprovider.ErrCodeNotAuthorizedException:
				return &cognitoidentityprovider.InitiateAuthOutput{}, http.StatusUnauthorized, errors.New("incorrect email or password")
			default:
				return &cognitoidentityprovider.InitiateAuthOutput{}, http.StatusInternalServerError, errors.New("something went wrong while signing in")
			}
		}
	}

	return authOutput, http.StatusOK, nil
}

func (c *Cognito) ConfirmSignUp(input *cognitoidentityprovider.ConfirmSignUpInput) (*cognitoidentityprovider.ConfirmSignUpOutput, int, error) {
	confirmOutput, err := c.Client.ConfirmSignUp(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cognitoidentityprovider.ErrCodeCodeMismatchException:
				return &cognitoidentityprovider.ConfirmSignUpOutput{}, http.StatusUnauthorized, errors.New("invalid verification code provided, please try again")
			case cognitoidentityprovider.ErrCodeExpiredCodeException:
				return &cognitoidentityprovider.ConfirmSignUpOutput{}, http.StatusUnauthorized, errors.New("verification code provided is expired, please try again from the start")
			default:
				return &cognitoidentityprovider.ConfirmSignUpOutput{}, http.StatusInternalServerError, errors.New("something went wrong during sign up")
			}
		}
	}

	return confirmOutput, http.StatusOK, nil
}

func (c *Cognito) ForgotPassword(input *cognitoidentityprovider.ForgotPasswordInput) (*cognitoidentityprovider.ForgotPasswordOutput, int, error) {
	forgotOutput, err := c.Client.ForgotPassword(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cognitoidentityprovider.ErrCodeUserNotFoundException:
				return &cognitoidentityprovider.ForgotPasswordOutput{}, http.StatusNotFound, errors.New("user not found")
			default:
				return &cognitoidentityprovider.ForgotPasswordOutput{}, http.StatusInternalServerError, errors.New("something went wrong")
			}
		}
	}

	return forgotOutput, http.StatusOK, nil
}

func (c *Cognito) GetUser(input *cognitoidentityprovider.GetUserInput) (*cognitoidentityprovider.GetUserOutput, int, error) {
	userOutput, err := c.Client.GetUser(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cognitoidentityprovider.ErrCodeUserNotFoundException:
				return &cognitoidentityprovider.GetUserOutput{}, http.StatusNotFound, errors.New("user not found")
			default:
				return &cognitoidentityprovider.GetUserOutput{}, http.StatusInternalServerError, errors.New("something went wrong while retrieving user")
			}
		}
	}

	return userOutput, http.StatusOK, nil
}

func (c *Cognito) ConfirmForgotPassword(input *cognitoidentityprovider.ConfirmForgotPasswordInput) (*cognitoidentityprovider.ConfirmForgotPasswordOutput, int, error) {
	confirmForgotOutput, err := c.Client.ConfirmForgotPassword(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cognitoidentityprovider.ErrCodeCodeMismatchException:
				return &cognitoidentityprovider.ConfirmForgotPasswordOutput{}, http.StatusUnauthorized, errors.New("invalid confirmation code")
			case cognitoidentityprovider.ErrCodeExpiredCodeException:
				return &cognitoidentityprovider.ConfirmForgotPasswordOutput{}, http.StatusUnauthorized, errors.New("confirmation code has expired")
			default:
				return &cognitoidentityprovider.ConfirmForgotPasswordOutput{}, http.StatusInternalServerError, errors.New("something went wrong while confirming password reset")
			}
		}
	}

	return confirmForgotOutput, http.StatusOK, nil
}
