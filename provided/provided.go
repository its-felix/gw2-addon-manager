package provided

import (
	"embed"
	"io"
)

//go:embed *.lua
var fs embed.FS

func LoadProvided(name string) (string, error) {
	f, err := fs.Open(name)
	if err != nil {
		return "", err
	}

	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
