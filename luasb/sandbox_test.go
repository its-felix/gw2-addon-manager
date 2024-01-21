package luasb

import (
	"context"
	"github.com/its-felix/shine"
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
	v, err := sb.Run(context.Background(), "return 123")

	if assert.NoError(t, err) {
		if v, ok := v.(lua.LNumber); ok {
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
		assert.Equal(t, v, lua.LNil)
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

	resCh := make(chan shine.Result[lua.LValue], numParallel)
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
		if assert.True(t, res.IsOk()) {
			v := res.Unwrap()
			if num, ok := res.Unwrap().(lua.LNumber); ok {
				assert.Equal(t, 1.0, float64(num))
			} else {
				assert.Failf(t, "expected return value to be a number", "was %v", v.Type())
			}
		}

		countRes++
	}

	assert.Equal(t, countRes, numParallel)
}

func TestRunFunc(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	sb := NewSandbox()
	v, err := sb.Run(ctx, "return function(arg) return arg end")

	if assert.NoError(t, err) {
		if fn, ok := v.(*lua.LFunction); ok {
			v, err = sb.RunFunc(ctx, fn, lua.LString("hello world"))

			if assert.NoError(t, err) {
				assert.Equal(t, lua.LString("hello world"), v)
			}
		} else {
			assert.Failf(t, "expected return value to be a function", "was %v", v.Type())
		}
	}
}

func assertIsTimeout(t testing.TB, err error) {
	var luaErr *lua.ApiError
	if assert.ErrorAs(t, err, &luaErr) {
		assert.Contains(t, err.Error(), "context deadline exceeded")
	}
}
