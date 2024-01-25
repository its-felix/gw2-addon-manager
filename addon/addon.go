package addon

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/its-felix/gw2-addon-manager/luasb"
	"github.com/its-felix/gw2-addon-manager/provided"
	lua "github.com/yuin/gopher-lua"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	gw2InstallDir  = `C:\Program Files (x86)\Steam\steamapps\common\Guild Wars 2`
	providedPrefix = "github.com/its-felix/gw2-addon-manager/provided/"
)

var (
	ErrInvalidScript = errors.New("invalid script")
)

type Addon struct {
	sb      *luasb.Sandbox
	version int
	name    string
	install *lua.LFunction
}

func (a *Addon) Version() int {
	return a.version
}

func (a *Addon) Name() string {
	return a.name
}

func (a *Addon) Install(ctx context.Context) error {
	_, err := a.sb.RunFunc(ctx, a.install, a.sb.Do(newApiTable))
	return err
}

func NewAddon(script string) (*Addon, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*5)
	defer cancel()

	sb := luasb.NewSandbox(withRequire())
	vs, err := sb.Run(ctx, script)
	if err != nil {
		return nil, err
	}

	if len(vs) != 1 {
		return nil, fmt.Errorf("invalid number of returns: %w", ErrInvalidScript)
	}

	tb, ok := vs[0].(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("the addon script returned a non-table value: %w", ErrInvalidScript)
	}

	a := &Addon{
		sb: sb,
	}

	lVersion, lName, lInstall := tb.RawGetString("version"), tb.RawGetString("name"), tb.RawGetString("install")
	if version, ok := lVersion.(lua.LNumber); ok {
		a.version = int(version)
	} else {
		return nil, fmt.Errorf("the addon script returned an invalid value for version: %w", ErrInvalidScript)
	}

	if name, ok := lName.(lua.LString); ok {
		a.name = string(name)
	} else {
		return nil, fmt.Errorf("the addon script returned an invalid value for name: %w", ErrInvalidScript)
	}

	if install, ok := lInstall.(*lua.LFunction); ok {
		a.install = install
	} else {
		return nil, fmt.Errorf("the addon script returned an invalid value for install: %w", ErrInvalidScript)
	}

	return a, nil
}

func withRequire() luasb.Option {
	return luasb.WithFunction("require", func(l *lua.LState) int {
		v := l.CheckString(1)
		name := strings.TrimPrefix(v, providedPrefix)
		if name == v {
			return 0
		}

		script, err := provided.LoadProvided(name)
		if err != nil {
			return 0
		}

		top := l.GetTop()
		if err = l.DoString(script); err != nil {
			return 0
		}

		numRet := l.GetTop() - top
		for i := 0; i < numRet; i++ {
			l.Push(l.Get(top + i + 1))
		}

		return numRet
	})
}

func newApiTable(l *lua.LState) lua.LValue {
	tb := l.NewTable()
	tb.RawSetString("download", l.NewFunction(func(l *lua.LState) int {
		url := l.CheckString(1)
		req, err := http.NewRequestWithContext(l.Context(), http.MethodGet, url, nil)
		if err != nil {
			l.RaiseError(err.Error())
			return 0
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			l.RaiseError(err.Error())
			return 0
		}

		defer res.Body.Close()

		b, err := io.ReadAll(res.Body)
		if err != nil {
			l.RaiseError(err.Error())
			return 0
		}

		ud := l.NewUserData()
		ud.Value = b

		tb := l.NewTable()
		tb.RawSetString("kind", lua.LString("httpresp"))
		tb.RawSetString("status", lua.LNumber(res.StatusCode))
		tb.RawSetString("body", ud)

		l.Push(tb)

		return 1
	}))

	tb.RawSetString("open", l.NewFunction(func(l *lua.LState) int {
		relPath := l.CheckString(1)
		path, err := filepath.Abs(filepath.Join(gw2InstallDir, relPath))
		if err != nil || !strings.HasPrefix(path, gw2InstallDir) {
			l.RaiseError("invalid path")
			return 0
		}

		f, err := os.Open(path)
		if err != nil {
			if os.IsNotExist(err) {
				l.Push(lua.LNil)
				return 1
			} else {
				l.RaiseError("invalid path")
				return 0
			}
		}

		defer f.Close()

		b, err := io.ReadAll(f)
		if err != nil {
			l.RaiseError(err.Error())
			return 0
		}

		ud := l.NewUserData()
		ud.Value = b

		tb := l.NewTable()
		tb.RawSetString("kind", lua.LString("file"))
		tb.RawSetString("body", ud)

		l.Push(tb)

		return 1
	}))

	tb.RawSetString("sha256", l.NewFunction(func(l *lua.LState) int {
		tb := l.CheckTable(1)
		lKind, ok := tb.RawGetString("kind").(lua.LString)
		if !ok {
			l.RaiseError("invalid param")
			return 0
		}

		var content []byte

		switch lKind {
		case "httpresp", "file":
			lBody, ok := tb.RawGetString("body").(*lua.LUserData)
			if !ok {
				l.RaiseError("invalid body")
				return 0
			}

			if content, ok = lBody.Value.([]byte); !ok {
				l.RaiseError("invalid body")
				return 0
			}

		default:
			l.RaiseError("invalid kind")
			return 0
		}

		hash := sha256.Sum256(content)
		l.Push(lua.LString(hex.EncodeToString(hash[:])))
		return 1
	}))

	return tb
}
