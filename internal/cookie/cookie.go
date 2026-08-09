package cookie

import (
	"net/http"
	"time"

	"github.com/spf13/viper"
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
	cookie.Expires = time.Now().UTC().Add(time.Hour)
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.Path = "/"
	cookie.Domain = viper.GetString("server.domain")
	cookie.SameSite = http.SameSiteLaxMode
	return cookie
}

func (cm *CookieManager) CookieClear() *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = CookieKey
	cookie.Value = ""
	cookie.Expires = time.Now().UTC().AddDate(-1, 0, 0)
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.Path = "/"
	// TODO: add cookie domain configuration
	cookie.Domain = viper.GetString("server.domain")
	cookie.SameSite = http.SameSiteLaxMode
	return cookie
}
