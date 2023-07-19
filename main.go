package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/hiepd/cognito-go"
)

const flowUsernamePassword = "USER_PASSWORD_AUTH"
const flowRefreshToken = "REFRESH_TOKEN_AUTH"

func main() {
	cognitoClientID := os.Getenv("COGNITO_CLIENT_ID")
	cognitoUserPoolID := os.Getenv("COGNITO_USER_POOL_ID")

	sessionCookie := Cookie{Name: "cognito-session"}
	tokenCookie := Cookie{Name: "cognito-token"}

	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	awsCogInstance := cognitoidentityprovider.New(sess)
	cognitoAuthClient, _ := cognito.NewCognitoClient(*sess.Config.Region, cognitoUserPoolID, cognitoClientID)

	r := gin.Default()
	r.LoadHTMLGlob("html/*")

	r.POST("/login", func(c *gin.Context) {
		flow := aws.String(flowUsernamePassword)
		params := map[string]*string{
			"USERNAME": aws.String(c.PostForm("username")),
			"PASSWORD": aws.String(c.PostForm("password")),
		}

		res, err := awsCogInstance.InitiateAuth(&cognitoidentityprovider.InitiateAuthInput{
			AuthFlow:       flow,
			AuthParameters: params,
			ClientId:       aws.String(cognitoClientID),
		})

		if err != nil {
			c.HTML(http.StatusOK, "login.html", gin.H{})
			return
		}

		if res.Session != nil {
			sessionCookie.Set(c, *res.Session)
		}

		if res.AuthenticationResult == nil {
			c.HTML(http.StatusOK, "changepassword.html", gin.H{})
			return
		}

		tokenCookie.Set(c, *res.AuthenticationResult.IdToken)
		c.Redirect(http.StatusSeeOther, "/home")
	})

	r.POST("/changepassword", func(c *gin.Context) {
		session, err := sessionCookie.Get(c)
		if err != nil {
			sessionCookie.Clear(c)
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}

		res, err := awsCogInstance.RespondToAuthChallenge(&cognitoidentityprovider.RespondToAuthChallengeInput{
			ChallengeName: aws.String("NEW_PASSWORD_REQUIRED"),
			ClientId:      aws.String(cognitoClientID),
			ChallengeResponses: map[string]*string{
				"USERNAME":     aws.String(c.PostForm("username")),
				"NEW_PASSWORD": aws.String(c.PostForm("new_password")),
			},
			Session: aws.String(session),
		})

		if err != nil {
			fmt.Println("Error on changepassword:", err)
			sessionCookie.Clear(c)
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}

		tokenCookie.Set(c, *res.AuthenticationResult.IdToken)
		c.Redirect(http.StatusSeeOther, "/home")
	})

	r.POST("/logout", func(c *gin.Context) {
		tokenCookie.Clear(c)
		c.Redirect(http.StatusSeeOther, "/login")
	})

	r.GET("/login", func(c *gin.Context) {
		token, err := tokenCookie.Get(c)
		if err == nil && token != "" {
			c.Redirect(http.StatusTemporaryRedirect, "/home")
			return
		}
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})

	r.GET("/home", func(c *gin.Context) {
		token, err := tokenCookie.Get(c)
		if err != nil {
			tokenCookie.Clear(c)
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}
		verifiedToken, err := cognitoAuthClient.VerifyToken(token)
		if err != nil || !verifiedToken.Valid {
			tokenCookie.Clear(c)
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}

		claims, ok := verifiedToken.Claims.(jwt.MapClaims)
		if !ok {
			tokenCookie.Clear(c)
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}
		username, ok := claims["cognito:username"].(string)
		if !ok {
			tokenCookie.Clear(c)
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}
		c.HTML(http.StatusOK, "home.html", gin.H{
			"username": username,
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
	})

	r.Run()
}
