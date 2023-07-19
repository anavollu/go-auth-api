package main

import (
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gin-gonic/gin"
)

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdatePassword struct {
	Username    string `json:"username"`
	NewPassword string `json:"newpassword"`
}

// TODO: Cookie on different file
type Cookie struct {
	Name string
}

func (ck Cookie) Set(c *gin.Context, token string) {
	cookie := http.Cookie{
		Name:     ck.Name,
		Value:    token,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Path:     "/",
	}
	http.SetCookie(c.Writer, &cookie)
}

func (ck Cookie) Get(c *gin.Context) (string, error) {
	cookie, err := c.Request.Cookie(ck.Name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (ck Cookie) Clear(c *gin.Context) {
	cookie := http.Cookie{
		Name:     ck.Name,
		Value:    "",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Path:     "/",
		MaxAge:   -1,
	}
	http.SetCookie(c.Writer, &cookie)
}

const flowUsernamePassword = "USER_PASSWORD_AUTH"
const flowRefreshToken = "REFRESH_TOKEN_AUTH"

func main() {
	cognitoClientID := os.Getenv("COGNITO_CLIENT_ID")

	sessionCookie := Cookie{Name: "cognito-session"}
	tokenCookie := Cookie{Name: "cognito-token"}

	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	cognito := cognitoidentityprovider.New(sess)

	r := gin.Default()
	r.LoadHTMLFiles("home.html")

	r.POST("/login", func(c *gin.Context) {
		loginReqBody := UserCredentials{}
		err := c.BindJSON(&loginReqBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": err.Error(),
			})
			return
		}

		flow := aws.String(flowUsernamePassword)
		params := map[string]*string{
			"USERNAME": aws.String(loginReqBody.Username),
			"PASSWORD": aws.String(loginReqBody.Password),
		}

		res, err := cognito.InitiateAuth(&cognitoidentityprovider.InitiateAuthInput{
			AuthFlow:       flow,
			AuthParameters: params,
			ClientId:       aws.String(cognitoClientID),
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": err.Error(),
			})
			return
		}

		if res.Session != nil {
			sessionCookie.Set(c, *res.Session)
		}

		if res.AuthenticationResult == nil {
			c.JSON(http.StatusForbidden, gin.H{
				"challenge": res.ChallengeName,
				"params":    res.ChallengeParameters,
			})
			return
		}

		tokenCookie.Set(c, *res.AuthenticationResult.IdToken)
		c.JSON(http.StatusOK, gin.H{
			"msg": "ok",
		})
	})

	r.POST("/changepassword", func(c *gin.Context) {
		reqBody := UpdatePassword{}
		err := c.BindJSON(&reqBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": err.Error(),
			})
			return
		}

		session, err := sessionCookie.Get(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": err.Error(),
			})
			return
		}

		res, err := cognito.RespondToAuthChallenge(&cognitoidentityprovider.RespondToAuthChallengeInput{
			ChallengeName: aws.String("NEW_PASSWORD_REQUIRED"),
			ClientId:      aws.String(cognitoClientID),
			ChallengeResponses: map[string]*string{
				"NEW_PASSWORD": aws.String(reqBody.NewPassword),
				"USERNAME":     aws.String(reqBody.Username),
			},
			Session: aws.String(session),
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": err.Error(),
			})
			return
		}

		tokenCookie.Set(c, *res.AuthenticationResult.IdToken)
		c.JSON(http.StatusOK, gin.H{
			"msg": "ok",
		})
	})

	r.GET("/home", func(c *gin.Context) {
		_, err := tokenCookie.Get(c)
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}

		// validate token

		// get the username from token

		c.HTML(http.StatusOK, "home.html", gin.H{
			"username": "John Doe", // Pass any required data to the template
		})
	})

	r.Run()
}
