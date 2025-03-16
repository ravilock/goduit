package cookie

import (
	"net/http"
	"time"
)

const CookieKey = "auth"

type CookieManager struct{}

func NewCookieManager() *CookieManager {
	return &CookieManager{}
}

func (cm *CookieManager) Create(token string) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = CookieKey
	cookie.Value = token
	cookie.Expires = time.Now().Add(time.Hour)
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.Path = "/"
	cookie.Domain = "localhost:3000"
	cookie.SameSite = http.SameSiteLaxMode
	return cookie
}

func (cm *CookieManager) CookieClear() *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = CookieKey
	cookie.Value = ""
	cookie.Expires = time.Now().AddDate(-1, 0, 0)
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.Path = "/"
	// TODO: add cookie domain configuration
	cookie.Domain = "localhost"
	cookie.SameSite = http.SameSiteLaxMode
	return cookie
}
