package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
