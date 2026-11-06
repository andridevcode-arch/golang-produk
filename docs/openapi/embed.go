package openapi

import (
	"embed"
	"fmt"
	"io/fs"
)

//go:embed *.yaml
var files embed.FS

// ReadFile returns one of the OpenAPI documents embedded in this package.
func ReadFile(name string) ([]byte, error) {
	data, err := fs.ReadFile(files, name)
	if err != nil {
		return nil, fmt.Errorf("read embedded OpenAPI document %q: %w", name, err)
	}
	return data, nil
}
