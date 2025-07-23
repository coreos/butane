package util

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func initGomplate(t *testing.T, gomplateConfig string) error {
	t.Helper()
	EnableGomplate = true
	once = sync.Once{} // reset for every test

	if gomplateConfig != "" {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, ".gomplate.yaml")

		err := os.WriteFile(configPath, []byte(gomplateConfig), 0644)
		if err != nil {
			return err
		}

		GomplateSourceDir = configPath
	} else {
		GomplateSourceDir = ""
	}

	return nil
}

func evalTemplate(t *testing.T, template string) (string, error) {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "template-*.tmpl")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(template); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	_, err = tmpFile.Seek(0, 0)
	if err != nil {
		t.Fatalf("failed to return to begining of temp file: %v", err)
	}

	output, err := GomplateReadLocalFile(tmpFile)
	return string(output), err
}

func TestInvalidGomplateConfig(t *testing.T) {
	if err := initGomplate(t, "not: valid: config: ["); err != nil {
		t.Fatalf("failed to write invalid config: %v", err)
	}

	renderer, err := GetGomplateRenderer()
	if err == nil {
		t.Fatalf("expected error for invalid config, got nil")
	}
	if renderer != nil {
		t.Fatalf("expected nil renderer, got non-nil")
	}
}

func TestInvalidTemplate(t *testing.T) {
	if err := initGomplate(t, "#empty config"); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	_, err := evalTemplate(t, "{{ .NonExistentField }}")

	if err == nil {
		t.Fatalf("expected error for missing key, got nil")
	}
}

func TestNoCustomConfig(t *testing.T) {
	if err := initGomplate(t, "#empty config"); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	rendered, err := evalTemplate(t, "{{ \"foobarbazquxquux\" | strings.Abbrev 9 }}")

	if err != nil {
		t.Fatalf("unexpected error: %+v\n", err)
	}

	if rendered != "foobar..." {
		t.Fatalf("Invalid rendered template, got: %s\n", rendered)
	}
}

func TestGomplateConfigApplication(t *testing.T) {
	// Create a mock HTTP server that returns a fixed JSON response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]string{"hello": "Hello"}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			t.Fatalf("json encoding failed: %v", err)
		}
	}))
	defer ts.Close()

	configContent := `
      leftDelim: ($(
      rightDelim: )$)
      context:
        data:
          url: ` + ts.URL
	if err := initGomplate(t, configContent); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	rendered, err := evalTemplate(t, "($( .data.hello )$)!")

	if err != nil {
		t.Fatalf("unexpected error: %+v\n", err)
	}
	if rendered != "Hello!" {
		t.Fatalf("Invalid rendered template, got: %s\n", rendered)
	}
}

func TestGomplateDisabled(t *testing.T) {
	EnableGomplate = false
	once = sync.Once{}

	expected := "some raw content"
	rendered, err := evalTemplate(t, expected)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rendered != expected {
		t.Fatalf("Invalid rendered template, got: %s\n", rendered)
	}
}

func TestMissingGomplateConfigFile(t *testing.T) {
	EnableGomplate = true
	once = sync.Once{}
	GomplateSourceDir = "/nonexistent/path/.gomplate.yaml"

	r, err := GetGomplateRenderer()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatalf("expected a valid renderer when config is missing")
	}
}
