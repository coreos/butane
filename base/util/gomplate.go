package util

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"sync"

	"github.com/hairyhenderson/gomplate/v4"
)

var once sync.Once

func initGomplateConfig() (*gomplate.Config, error) {
	f, err := os.Open(".gomplate.yaml")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	return gomplate.Parse(f)
}

func createGomplateRenderer() (gomplate.Renderer, error) {
	_, err := initGomplateConfig()

	if err != nil {
		return nil, err
	}

	return gomplate.NewRenderer(gomplate.RenderOptions{}), nil
}

func GetGomplateRenderer() (gomplate.Renderer, error) {
	var renderer gomplate.Renderer
	var err error
	once.Do(func() {
		renderer, err = createGomplateRenderer()
	})
	return renderer, err
}

func GomplateReadLocalFile(file *os.File) ([]byte, error) {
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	gomplateRenderer, _ := GetGomplateRenderer()
	// gomplate is disabled
	if gomplateRenderer == nil {
		return fileContent, nil
	}

	var buf bytes.Buffer
	err = gomplateRenderer.Render(context.Background(), file.Name(), string(fileContent), &buf)
	return buf.Bytes(), err
}
