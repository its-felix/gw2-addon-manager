package web

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

var csp = strings.Join([]string{
	"default-src 'self'",
	"connect-src 'self'",
	"style-src 'self' 'unsafe-inline'",
	"font-src data:",
	"img-src 'self' https://static.staticwars.com/quaggans/",
	"script-src 'self' 'unsafe-inline'",
}, "; ")

func HeadersMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set(echo.HeaderCacheControl, "private, no-cache, no-store, max-age=0, must-revalidate")
			c.Response().Header().Set(echo.HeaderContentSecurityPolicy, csp)

			return next(c)
		}
	}
}

func AuthenticatedMiddleware(cookieName, token string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(cookieName)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err)
			}

			if err = cookie.Valid(); err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err)
			}

			if cookie.Value != token {
				return echo.NewHTTPError(http.StatusUnauthorized, err)
			}

			return next(c)
		}
	}
}
