package common

import (
	"net/http"
	"time"
)

// func SetCookie(
//
//	key, value string,
//	exp_mins int,
//
//	) *http.Cookie {
//		cookie := new(http.Cookie)
//		cookie.Name = key
//		cookie.Value = value
//		cookie.Path = "/"
//		cookie.HttpOnly = true
//		cookie.Secure = false
//		cookie.SameSite = http.SameSiteNoneMode
//		cookie.Expires = time.Now().Add(time.Minute * time.Duration(exp_mins))
//		return cookie
//	}

func SetCookie(key, value string, exp_mins int) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = key
	cookie.Value = value
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = false
	cookie.SameSite = http.SameSiteLaxMode
	cookie.Expires = time.Now().Add(time.Minute * time.Duration(exp_mins))
	return cookie
}
