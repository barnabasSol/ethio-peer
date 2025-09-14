package common

import (
	"net/http"
	"time"
)

func SetCookie(key, value string, exp_mins int) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = key
	cookie.Value = value
	cookie.Path = "/"
	cookie.HttpOnly = true
	// cookie.Domain = ".barney-host.site"
	cookie.Secure = false
	cookie.SameSite = http.SameSiteLaxMode
	cookie.Expires = time.Now().Add(time.Minute * time.Duration(exp_mins))
	return cookie
}
