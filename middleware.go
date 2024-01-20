package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

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
