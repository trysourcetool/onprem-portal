package server

import (
	"net/http"

	"github.com/trysourcetool/onprem-portal/internal/config"
)

const domain = "portal.trysourcetool.com"

type cookieConfig struct {
	isLocalEnv bool
}

func newCookieConfig() *cookieConfig {
	return &cookieConfig{
		isLocalEnv: config.Config.Env == config.EnvLocal,
	}
}

func (c *cookieConfig) getXSRFTokenSameSite() http.SameSite {
	if c.isLocalEnv {
		return http.SameSiteLaxMode
	}
	return http.SameSiteNoneMode
}

func (c *cookieConfig) isSecure() bool {
	return !c.isLocalEnv
}

func (c *cookieConfig) setCookie(w http.ResponseWriter, name, value string, maxAge int, httpOnly bool, sameSite http.SameSite) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		Path:     "/",
		Domain:   domain,
		HttpOnly: httpOnly,
		Secure:   c.isSecure(),
		SameSite: sameSite,
	})
}

func (c *cookieConfig) deleteCookie(w http.ResponseWriter, r *http.Request, name string, httpOnly bool, sameSite http.SameSite) {
	if cookie, _ := r.Cookie(name); cookie != nil {
		cookie.MaxAge = -1
		cookie.Domain = domain
		cookie.Path = "/"
		cookie.HttpOnly = httpOnly
		cookie.Secure = c.isSecure()
		cookie.SameSite = sameSite
		http.SetCookie(w, cookie)
	}
}

func (c *cookieConfig) SetAuthCookie(w http.ResponseWriter, token, refreshToken, xsrfToken string, tokenMaxAge, refreshTokenMaxAge, xsrfTokenMaxAge int) {
	xsrfTokenSameSite := c.getXSRFTokenSameSite()

	c.setCookie(w, "access_token", token, tokenMaxAge, true, http.SameSiteStrictMode)
	c.setCookie(w, "refresh_token", refreshToken, refreshTokenMaxAge, true, http.SameSiteStrictMode)
	c.setCookie(w, "xsrf_token", xsrfToken, xsrfTokenMaxAge, false, xsrfTokenSameSite)
	c.setCookie(w, "xsrf_token_same_site", xsrfToken, xsrfTokenMaxAge, true, http.SameSiteStrictMode)
}

func (c *cookieConfig) DeleteAuthCookie(w http.ResponseWriter, r *http.Request) {
	xsrfTokenSameSite := c.getXSRFTokenSameSite()

	c.deleteCookie(w, r, "access_token", true, http.SameSiteStrictMode)
	c.deleteCookie(w, r, "refresh_token", true, http.SameSiteStrictMode)
	c.deleteCookie(w, r, "xsrf_token", false, xsrfTokenSameSite)
	c.deleteCookie(w, r, "xsrf_token_same_site", true, http.SameSiteStrictMode)
}
