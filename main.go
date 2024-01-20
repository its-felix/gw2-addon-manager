package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
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

	e, token, err := setup(stop)
	if err != nil {
		panic(err)
	}

	errCh := make(chan error)
	go func() {
		defer close(errCh)

		if err := run(ctx, e); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	if err = waitUntilHealthy(ctx, time.Second*5); err != nil {
		panic(err)
	}

	if err = open(fmt.Sprintf("http://%s/init/%s", listenAddr, url.PathEscape(token))); err != nil {
		panic(err)
	}

	if err = <-errCh; err != nil {
		panic(err)
	}
}

func setup(shutdownFunc func()) (*echo.Echo, string, error) {
	proxyUrl, err := url.Parse("http://127.0.0.1:4200/")
	if err != nil {
		return nil, "", err
	}

	token, err := generateToken()
	if err != nil {
		return nil, "", err
	}

	e := echo.New()
	e.Group("/", middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{{URL: proxyUrl}})))

	e.HEAD("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	e.GET("/init/:token", func(c echo.Context) error {
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

	apiGroup := e.Group("/api", AuthenticatedMiddleware(tokenCookieName, token))
	apiGroup.POST("/shutdown", func(c echo.Context) error {
		defer shutdownFunc()
		return c.NoContent(http.StatusOK)
	})

	return e, token, nil
}

func generateToken() (string, error) {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func run(ctx context.Context, e *echo.Echo) error {
	done := make(chan error)
	go func() {
		defer close(done)

		<-ctx.Done()
		if err := e.Shutdown(context.Background()); err != nil {
			done <- err
		}
	}()

	if err := e.Start(listenAddr); err != nil {
		return err
	}

	return <-done
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

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, fmt.Sprintf("http://%s/health", listenAddr), nil)
	if err != nil {
		return err
	}

	var res *http.Response
	for {
		res, _ = http.DefaultClient.Do(req)
		if res.StatusCode == http.StatusOK {
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
