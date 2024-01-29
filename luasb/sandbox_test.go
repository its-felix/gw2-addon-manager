package luasb

import (
	"context"
	"github.com/its-felix/shine/v2"
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	sb := NewSandbox()
	_, err := sb.Run(ctx, "while true do end")
	assertIsTimeout(t, err)
}

func TestReturn(t *testing.T) {
	sb := NewSandbox()
	vs, err := sb.Run(context.Background(), "return 123")

	if assert.NoError(t, err) && assert.Len(t, vs, 1) {
		if v, ok := vs[0].(lua.LNumber); ok {
			assert.Equal(t, float64(v), 123.0)
		} else {
			assert.FailNow(t, "expected return value to be a number")
		}
	}
}

func TestNoReturn(t *testing.T) {
	sb := NewSandbox()
	v, err := sb.Run(context.Background(), "")

	if assert.NoError(t, err) {
		assert.Len(t, v, 0)
	}
}

func TestParallel(t *testing.T) {
	var running atomic.Uint32
	var wgStarted sync.WaitGroup
	var wgDone sync.WaitGroup

	numParallel := 10
	wgStarted.Add(numParallel)

	sb := NewSandbox(
		WithFunction("run", func(l *lua.LState) int {
			wgStarted.Wait()

			v := running.Add(1)
			defer running.Add(^uint32(0))

			time.Sleep(time.Millisecond * 10)
			l.Push(lua.LNumber(v))
			return 1
		}),
	)

	resCh := make(chan shine.Result[[]lua.LValue], numParallel)
	for i := 0; i < numParallel; i++ {
		wgDone.Add(1)
		go func() {
			defer wgDone.Done()
			wgStarted.Done()
			resCh <- shine.NewResult(sb.Run(context.Background(), "return run()"))
		}()
	}

	wgDone.Wait()
	close(resCh)

	countRes := 0
	for res := range resCh {
		if vs, _, ok := res.Get(); assert.True(t, ok) {
			if assert.Len(t, vs, 1) {
				if num, ok := vs[0].(lua.LNumber); ok {
					assert.Equal(t, 1.0, float64(num))
				} else {
					assert.Failf(t, "expected return value to be a number", "was %v", vs[0].Type())
				}
			}
		}

		countRes++
	}

	assert.Equal(t, numParallel, countRes)
}

func TestRunFunc(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	sb := NewSandbox()
	vs, err := sb.Run(ctx, "return function(...) return ... end")

	if assert.NoError(t, err) && assert.Len(t, vs, 1) {
		if fn, ok := vs[0].(*lua.LFunction); ok {
			vs, err = sb.RunFunc(ctx, fn, lua.LString("hello"), lua.LString("world"))

			if assert.NoError(t, err) && assert.Len(t, vs, 2) {
				assert.Equal(t, lua.LString("hello"), vs[0])
				assert.Equal(t, lua.LString("world"), vs[1])
			}
		} else {
			assert.Failf(t, "expected return value to be a function", "was %v", vs[0].Type())
		}
	}
}

func TestVararg(t *testing.T) {
	sb := NewSandbox()
	vs, err := sb.Run(context.Background(), "return 1,2,3")

	if assert.NoError(t, err) && assert.Len(t, vs, 3) {
		assert.Equal(t, lua.LNumber(1), vs[0])
		assert.Equal(t, lua.LNumber(2), vs[1])
		assert.Equal(t, lua.LNumber(3), vs[2])
	}
}

func assertIsTimeout(t testing.TB, err error) {
	var luaErr *lua.ApiError
	if assert.ErrorAs(t, err, &luaErr) {
		assert.Contains(t, err.Error(), "context deadline exceeded")
	}
}
