package addon

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewTimeout(t *testing.T) {
	_, err := NewAddon("while true do end")
	assert.Error(t, err)
}

func TestNewMultipleReturns(t *testing.T) {
	_, err := NewAddon("return 1,2")
	assert.ErrorIs(t, err, ErrInvalidScript)
}

func TestNewNoVersion(t *testing.T) {
	_, err := NewAddon("return {}")
	assert.ErrorIs(t, err, ErrInvalidScript)
}

func TestNewNoName(t *testing.T) {
	_, err := NewAddon("return {version = 1}")
	assert.ErrorIs(t, err, ErrInvalidScript)
}

func TestNewAddon(t *testing.T) {
	a, err := NewAddon(`return {version = 1, name = "Test", install = function() end}`)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, a.version)
		assert.Equal(t, "Test", a.name)
	}
}

func TestRequire(t *testing.T) {
	const script = `
local base = require("github.com/its-felix/gw2-addon-manager/provided/test.lua")
return {
	version = base.version + 1,
	name = base.name .. " Test",
	install = base.install,
}
`

	a, err := NewAddon(script)
	if assert.NoError(t, err) {
		assert.Equal(t, 2, a.version)
		assert.Equal(t, "Test Test", a.name)
	}
}

func TestRequireTimeout(t *testing.T) {
	_, err := NewAddon(`require("github.com/its-felix/gw2-addon-manager/provided/test_timeout.lua")`)
	assert.Error(t, err)
}

func TestInstall(t *testing.T) {
	const script = `
return {
	version = 1,
	name = "Test",
	install = function(api)
		print(api.sha256(api.open("addons/arcdps/drf.dll")))
		print(api.sha256(api.download("https://update.drf.rs/drf.dll")))
	end,
}
`
	a, err := NewAddon(script)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, a.version)
		assert.Equal(t, "Test", a.name)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		assert.NoError(t, a.Install(ctx))
	}
}
