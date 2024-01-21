package web

import (
	"context"
	"github.com/labstack/echo/v4"
)

type Server struct {
	e *echo.Echo
}

func (s *Server) Run(ctx context.Context, addr string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	done := make(chan error)
	go func() {
		defer close(done)

		<-ctx.Done()
		if err := s.e.Shutdown(context.Background()); err != nil {
			done <- err
		}
	}()

	if err := s.e.Start(addr); err != nil {
		return err
	}

	return <-done
}

func NewServer(e *echo.Echo) (*Server, error) {
	return &Server{e}, nil
}
