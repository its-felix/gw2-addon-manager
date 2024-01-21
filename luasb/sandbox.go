package luasb

import (
	"context"
	lua "github.com/yuin/gopher-lua"
	"sync"
)

type Sandbox struct {
	l     *lua.LState
	mutex sync.Mutex
}

func (s *Sandbox) Run(ctx context.Context, script string) ([]lua.LValue, error) {
	return s.run(ctx, func() error {
		return s.l.DoString(script)
	})
}

func (s *Sandbox) RunFunc(ctx context.Context, fn *lua.LFunction, args ...lua.LValue) ([]lua.LValue, error) {
	return s.run(ctx, func() error {
		s.l.Push(fn)

		for _, arg := range args {
			s.l.Push(arg)
		}

		s.l.Call(len(args), lua.MultRet)
		return nil
	})
}

func (s *Sandbox) run(ctx context.Context, fn func() error) ([]lua.LValue, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.l.SetContext(ctx)
	defer s.l.RemoveContext()

	top := s.l.GetTop()
	if err := fn(); err != nil {
		return nil, err
	}

	numRet := s.l.GetTop() - top
	res := make([]lua.LValue, numRet)

	for i := 0; i < numRet; i++ {
		res[i] = s.l.Get(top + i + 1)
	}

	return res, nil
}

type Option func(l *lua.LState)

func NewSandbox(options ...Option) *Sandbox {
	l := lua.NewState()
	l.SetGlobal("io", lua.LNil)

	for _, option := range options {
		option(l)
	}

	return &Sandbox{
		l: l,
	}
}

func WithFunction(name string, fn lua.LGFunction) Option {
	return func(l *lua.LState) {
		l.SetGlobal(name, l.NewFunction(fn))
	}
}
