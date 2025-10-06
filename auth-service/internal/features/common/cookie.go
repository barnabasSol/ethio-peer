package common

import (
	"net/http"
	"time"
)

// func SetCookie(
// 	key, value string,
// 	exp_mins int,
// ) *http.Cookie {
// 	cookie := new(http.Cookie)
// 	cookie.Name = key
// 	cookie.Value = value
// 	cookie.Path = "/"
// 	cookie.Domain = "barney-host.site"
// 	cookie.HttpOnly = true
// 	cookie.Secure = true
// 	cookie.SameSite = http.SameSiteNoneMode
// 	cookie.Expires = time.Now().Add(time.Minute * time.Duration(exp_mins))
// 	return cookie
// }

func SetCookie(
	key, value string,
	exp_mins int,
) *http.Cookie {
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

func ClearCookie(key string) *http.Cookie {
	return &http.Cookie{
		Name:     key,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
	}
}
