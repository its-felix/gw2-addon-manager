package web

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

const CSP = "default-src 'self'"

func HeadersMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set(echo.HeaderCacheControl, "private, no-cache, no-store, max-age=0, must-revalidate")
			c.Response().Header().Set(echo.HeaderContentSecurityPolicy, CSP)

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
