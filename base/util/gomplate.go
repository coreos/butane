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
var EnableGomplate = false
var GomplateSourceDir string
var renderer gomplate.Renderer
var err error

func initGomplateConfig() (*gomplate.Config, error) {
	f, err := os.Open(GomplateSourceDir)
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
	if config, err := initGomplateConfig(); err != nil {
		return nil, err
	} else if config != nil {
		// TODO: is there a better way?
		return gomplate.NewRenderer(gomplate.RenderOptions{
			Datasources:  config.DataSources,
			Context:      config.Context,
			Templates:    config.Templates,
			ExtraHeaders: config.ExtraHeaders,
			LDelim:       config.LDelim,
			RDelim:       config.RDelim,
			MissingKey:   config.MissingKey,
		}), nil
	} else {
		return gomplate.NewRenderer(gomplate.RenderOptions{}), nil
	}
}

func GetGomplateRenderer() (gomplate.Renderer, error) {
	if !EnableGomplate {
		return nil, nil
	}

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

	// Initialization has been done during argument parsing
	// Any error would've happend then, thus we can safely ignore errors
	gomplateRenderer, _ := GetGomplateRenderer()
	// gomplate is disabled
	if gomplateRenderer == nil {
		return fileContent, nil
	}

	var buf bytes.Buffer
	err = gomplateRenderer.Render(context.Background(), file.Name(), string(fileContent), &buf)
	return buf.Bytes(), err
}
