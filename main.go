package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/its-felix/gw2-addon-manager/web"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"time"
)

const (
	listenAddr      = "127.0.0.1:8080"
	tokenCookieName = "token"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if isHealthy(ctx) {
		if err := open(fmt.Sprintf("http://%s/", listenAddr)); err != nil {
			panic(err)
		}
		return
	}

	s, token, err := setup(stop)
	if err != nil {
		panic(err)
	}

	errCh := make(chan error)
	go func() {
		defer close(errCh)

		if err := s.Run(ctx, listenAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	if err = waitUntilHealthy(ctx, time.Second*5); err != nil {
		panic(err)
	}

	if err = open(fmt.Sprintf("http://%s/auth/%s", listenAddr, url.PathEscape(token))); err != nil {
		panic(err)
	}

	if err = <-errCh; err != nil {
		panic(err)
	}
}

func setup(shutdownFn func()) (*web.Server, string, error) {
	proxyUrl, err := url.Parse("http://127.0.0.1:4200/")
	if err != nil {
		return nil, "", err
	}

	token, err := generateToken()
	if err != nil {
		return nil, "", err
	}

	e := echo.New()
	e.Use(web.HeadersMiddleware())

	e.Group("/", middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{{URL: proxyUrl}})))

	e.HEAD("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	e.GET("/auth/:token", func(c echo.Context) error {
		reqToken := c.Param("token")
		if reqToken != token {
			c.SetCookie(&http.Cookie{
				Name:   tokenCookieName,
				Value:  "",
				Path:   "/",
				MaxAge: -1,
			})

			return c.Redirect(http.StatusFound, "/error")
		}

		c.SetCookie(&http.Cookie{
			Name:     tokenCookieName,
			Value:    reqToken,
			Path:     "/",
			Expires:  time.Now().Add(time.Hour * 24),
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})

		return c.Redirect(http.StatusFound, "/")
	})

	apiGroup := e.Group("/api", web.AuthenticatedMiddleware(tokenCookieName, token))
	apiGroup.POST("/shutdown", func(c echo.Context) error {
		defer shutdownFn()
		return c.NoContent(http.StatusOK)
	})

	s, err := web.NewServer(e)
	if err != nil {
		return nil, "", err
	}

	return s, token, nil
}

func generateToken() (string, error) {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}

	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func waitUntilHealthy(ctx context.Context, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for {
		if isHealthy(ctx) {
			return nil
		}

		select {
		case <-time.After(time.Millisecond * 50):
			continue

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func isHealthy(ctx context.Context) bool {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, fmt.Sprintf("http://%s/health", listenAddr), nil)
	if err != nil {
		return false
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}

	return res.StatusCode == http.StatusOK
}
