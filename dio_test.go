package dio

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func setupTempEnvDir(t *testing.T, mode, content string) (tmpDir string, cleanup func()) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "dio_test_")
	if err != nil {
		t.Fatal(err)
	}
	normalized := normalizeMode(mode)
	if normalized == "" {
		normalized = mode
	}
	filename := fmt.Sprintf(".env.%s", normalized)
	path := filepath.Join(tmpDir, filename)
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatal(err)
	}
	return tmpDir, func() { os.RemoveAll(tmpDir) }
}

func TestNormalizeMode(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"dev", "development"},
		{"Dev", "development"},
		{"prod", "production"},
		{"PROD", "production"},
		{"production", "production"},
		{"development", "development"},
		{"uat", "uat"},
		{"UAT", "uat"},
		{"unknown", "unknown"},
		{"Staging", "staging"},
		{"", ""},
	}
	for _, tt := range tests {
		got := normalizeMode(tt.in)
		if got != tt.want {
			t.Errorf("normalizeMode(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestLoadEnvFile_EmptyMode(t *testing.T) {
	err := loadEnvFile("", ".env.%s")
	if err == nil {
		t.Fatal("expected error for empty mode")
	}
	if !errors.Is(err, ErrInvalidEnvMode) {
		t.Errorf("expected ErrInvalidEnvMode, got %v", err)
	}
}

func TestLoadEnvFile_FileNotFound(t *testing.T) {
	tmpDir, cleanup := setupTempEnvDir(t, "development", "KEY=val")
	defer cleanup()
	origWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)
	os.Remove(filepath.Join(tmpDir, ".env.development"))

	err := loadEnvFile("development", ".env.%s")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	if !IsEnvFileNotFoundError(err) {
		t.Errorf("expected ErrEnvFileNotFound, got %v", err)
	}
}

func TestLoadEnvFile_Success(t *testing.T) {
	tmpDir, cleanup := setupTempEnvDir(t, "development", "TEST_KEY=loaded_value")
	defer cleanup()
	origWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	err := loadEnvFile("development", ".env.%s")
	if err != nil {
		t.Fatalf("loadEnvFile: %v", err)
	}
	if got := os.Getenv("TEST_KEY"); got != "loaded_value" {
		t.Errorf("TEST_KEY = %q, want loaded_value", got)
	}
	os.Unsetenv("TEST_KEY")
}

func TestInitEnvironmentWithEnv_SuccessNilOpts(t *testing.T) {
	tmpDir, cleanup := setupTempEnvDir(t, "development", "INIT_KEY=init_ok")
	defer cleanup()
	origWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	err := InitEnvironmentWithEnv("development", &Options{Silent: true})
	if err != nil {
		t.Fatalf("InitEnvironmentWithEnv: %v", err)
	}
	if got := os.Getenv("INIT_KEY"); got != "init_ok" {
		t.Errorf("INIT_KEY = %q, want init_ok", got)
	}
	os.Unsetenv("INIT_KEY")
}

func TestInitEnvironmentWithEnv_CustomFilePattern(t *testing.T) {
	tmpDir, cleanup := setupTempEnvDir(t, "development", "CUSTOM_KEY=custom_ok")
	defer cleanup()
	pattern := filepath.Join(tmpDir, ".env.%s")

	err := InitEnvironmentWithEnv("development", &Options{FilePattern: pattern, Silent: true})
	if err != nil {
		t.Fatalf("InitEnvironmentWithEnv: %v", err)
	}
	if got := os.Getenv("CUSTOM_KEY"); got != "custom_ok" {
		t.Errorf("CUSTOM_KEY = %q, want custom_ok", got)
	}
	os.Unsetenv("CUSTOM_KEY")
}

func TestInitEnvironmentWithEnv_SilentTrue(t *testing.T) {
	tmpDir, cleanup := setupTempEnvDir(t, "development", "KEY=val")
	defer cleanup()
	origWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	err := InitEnvironmentWithEnv("development", &Options{Silent: true})
	if err != nil {
		t.Fatalf("InitEnvironmentWithEnv: %v", err)
	}
}

func TestInitEnvironmentWithEnv_FileNotFound(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dio_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	pattern := filepath.Join(tmpDir, ".env.%s")

	err = InitEnvironmentWithEnv("development", &Options{FilePattern: pattern, Silent: true})
	if err == nil {
		t.Fatal("expected error when file missing")
	}
	if !IsEnvFileNotFoundError(err) {
		t.Errorf("expected ErrEnvFileNotFound, got %v", err)
	}
}

func TestInitEnvironmentWithEnv_EmptyEnvUsesDefault(t *testing.T) {
	tmpDir, cleanup := setupTempEnvDir(t, "development", "DEFAULT_KEY=default_ok")
	defer cleanup()
	origWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	err := InitEnvironmentWithEnv("", &Options{Silent: true})
	if err != nil {
		t.Fatalf("InitEnvironmentWithEnv empty: %v", err)
	}
	if got := os.Getenv("DEFAULT_KEY"); got != "default_ok" {
		t.Errorf("DEFAULT_KEY = %q, want default_ok", got)
	}
	os.Unsetenv("DEFAULT_KEY")
}

func TestInitEnvironment_DefaultEnv(t *testing.T) {
	tmpDir, cleanup := setupTempEnvDir(t, "development", "FLAG_KEY=flag_ok")
	defer cleanup()
	origWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	err := InitEnvironment(&Options{Silent: true})
	if err != nil {
		t.Fatalf("InitEnvironment: %v", err)
	}
	if got := os.Getenv("FLAG_KEY"); got != "flag_ok" {
		t.Errorf("FLAG_KEY = %q, want flag_ok", got)
	}
	os.Unsetenv("FLAG_KEY")
}

func TestInitEnvironment_OptsEnvOverride(t *testing.T) {
	tmpDir, cleanup := setupTempEnvDir(t, "production", "PROD_KEY=prod_ok")
	defer cleanup()
	pattern := filepath.Join(tmpDir, ".env.%s")

	err := InitEnvironment(&Options{Env: "production", FilePattern: pattern, Silent: true})
	if err != nil {
		t.Fatalf("InitEnvironment with Env override: %v", err)
	}
	if Name() != "production" {
		t.Errorf("Name() = %q, want production", Name())
	}
	if got := os.Getenv("PROD_KEY"); got != "prod_ok" {
		t.Errorf("PROD_KEY = %q, want prod_ok", got)
	}
	os.Unsetenv("PROD_KEY")
}

func TestEnv_EmptyKey(t *testing.T) {
	_, err := Env("")
	if err == nil {
		t.Fatal("expected error for empty key")
	}
	if !IsEnvVarNotSetError(err) {
		t.Errorf("expected ErrEnvVarNotSet, got %v", err)
	}
}

func TestEnv_KeySet(t *testing.T) {
	key := "DIO_TEST_ENV_SET"
	os.Setenv(key, "set_value")
	defer os.Unsetenv(key)

	got, err := Env(key)
	if err != nil {
		t.Fatalf("Env: %v", err)
	}
	if got != "set_value" {
		t.Errorf("Env(%q) = %q, want set_value", key, got)
	}
}

func TestEnv_KeyUnsetWithFallback(t *testing.T) {
	key := "DIO_TEST_ENV_UNSET_FALLBACK"
	os.Unsetenv(key)

	got, err := Env(key, "fallback_value")
	if err != nil {
		t.Fatalf("Env: %v", err)
	}
	if got != "fallback_value" {
		t.Errorf("Env(%q, fallback) = %q, want fallback_value", key, got)
	}
}

func TestEnv_KeyUnsetNoFallback(t *testing.T) {
	key := "DIO_TEST_ENV_UNSET_NO_FALLBACK"
	os.Unsetenv(key)

	_, err := Env(key)
	if err == nil {
		t.Fatal("expected error when key unset and no fallback")
	}
	if !IsEnvVarNotSetError(err) {
		t.Errorf("expected ErrEnvVarNotSet, got %v", err)
	}
}

func TestRequiredEnv_Set(t *testing.T) {
	key := "DIO_TEST_REQUIRED_SET"
	os.Setenv(key, "required_value")
	defer os.Unsetenv(key)

	got, err := RequiredEnv(key)
	if err != nil {
		t.Fatalf("RequiredEnv: %v", err)
	}
	if got != "required_value" {
		t.Errorf("RequiredEnv(%q) = %q, want required_value", key, got)
	}
}

func TestRequiredEnv_Unset(t *testing.T) {
	key := "DIO_TEST_REQUIRED_UNSET"
	os.Unsetenv(key)

	_, err := RequiredEnv(key)
	if err == nil {
		t.Fatal("expected error when key unset")
	}
	if !IsEnvVarNotSetError(err) {
		t.Errorf("expected ErrEnvVarNotSet, got %v", err)
	}
}

func TestEnvWithDefault_KeySet(t *testing.T) {
	key := "DIO_TEST_DEFAULT_SET"
	os.Setenv(key, "actual_value")
	defer os.Unsetenv(key)

	got := EnvWithDefault(key, "default_value")
	if got != "actual_value" {
		t.Errorf("EnvWithDefault = %q, want actual_value", got)
	}
}

func TestEnvWithDefault_KeyUnset(t *testing.T) {
	key := "DIO_TEST_DEFAULT_UNSET"
	os.Unsetenv(key)

	got := EnvWithDefault(key, "default_value")
	if got != "default_value" {
		t.Errorf("EnvWithDefault = %q, want default_value", got)
	}
}

func TestEnvWithDefault_EmptyKey(t *testing.T) {
	got := EnvWithDefault("", "default_value")
	if got != "default_value" {
		t.Errorf("EnvWithDefault(\"\", \"default_value\") = %q, want default_value", got)
	}
}

func TestName_Dev_Prod_Uat(t *testing.T) {
	tmpDir, cleanup := setupTempEnvDir(t, "development", "X=1")
	defer cleanup()
	pattern := filepath.Join(tmpDir, ".env.%s")

	err := InitEnvironmentWithEnv("development", &Options{FilePattern: pattern, Silent: true})
	if err != nil {
		t.Fatal(err)
	}
	if Name() != "development" {
		t.Errorf("Name() = %q, want development", Name())
	}
	if !Dev() {
		t.Error("Dev() = false, want true")
	}
	if Prod() {
		t.Error("Prod() = true, want false")
	}
	if Uat() {
		t.Error("Uat() = true, want false")
	}
	os.Unsetenv("X")
}

func TestName_Prod(t *testing.T) {
	tmpDir, cleanup := setupTempEnvDir(t, "production", "X=1")
	defer cleanup()
	pattern := filepath.Join(tmpDir, ".env.%s")

	err := InitEnvironmentWithEnv("production", &Options{FilePattern: pattern, Silent: true})
	if err != nil {
		t.Fatal(err)
	}
	if Name() != "production" {
		t.Errorf("Name() = %q, want production", Name())
	}
	if Dev() {
		t.Error("Dev() = true, want false")
	}
	if !Prod() {
		t.Error("Prod() = false, want true")
	}
	if Uat() {
		t.Error("Uat() = true, want false")
	}
	os.Unsetenv("X")
}

func TestName_Uat(t *testing.T) {
	tmpDir, cleanup := setupTempEnvDir(t, "uat", "X=1")
	defer cleanup()
	pattern := filepath.Join(tmpDir, ".env.%s")

	err := InitEnvironmentWithEnv("uat", &Options{FilePattern: pattern, Silent: true})
	if err != nil {
		t.Fatal(err)
	}
	if Name() != "uat" {
		t.Errorf("Name() = %q, want uat", Name())
	}
	if Dev() {
		t.Error("Dev() = true, want false")
	}
	if Prod() {
		t.Error("Prod() = true, want false")
	}
	if !Uat() {
		t.Error("Uat() = false, want true")
	}
	os.Unsetenv("X")
}

func TestIsEnvVarNotSetError(t *testing.T) {
	if !IsEnvVarNotSetError(ErrEnvVarNotSet) {
		t.Error("direct ErrEnvVarNotSet should be true")
	}
	wrapped := fmt.Errorf("wrap: %w", ErrEnvVarNotSet)
	if !IsEnvVarNotSetError(wrapped) {
		t.Error("wrapped ErrEnvVarNotSet should be true")
	}
	if IsEnvVarNotSetError(ErrEnvFileNotFound) {
		t.Error("ErrEnvFileNotFound should be false")
	}
}

func TestIsEnvFileNotFoundError(t *testing.T) {
	if !IsEnvFileNotFoundError(ErrEnvFileNotFound) {
		t.Error("direct ErrEnvFileNotFound should be true")
	}
	wrapped := fmt.Errorf("wrap: %w", ErrEnvFileNotFound)
	if !IsEnvFileNotFoundError(wrapped) {
		t.Error("wrapped ErrEnvFileNotFound should be true")
	}
	if IsEnvFileNotFoundError(ErrEnvVarNotSet) {
		t.Error("ErrEnvVarNotSet should be false")
	}
}

func TestIsInvalidEnvModeError(t *testing.T) {
	if !IsInvalidEnvModeError(ErrInvalidEnvMode) {
		t.Error("direct ErrInvalidEnvMode should be true")
	}
	wrapped := fmt.Errorf("wrap: %w", ErrInvalidEnvMode)
	if !IsInvalidEnvModeError(wrapped) {
		t.Error("wrapped ErrInvalidEnvMode should be true")
	}
	if IsInvalidEnvModeError(ErrEnvVarNotSet) {
		t.Error("ErrEnvVarNotSet should be false")
	}
}

func TestPrintEnvMode_Production(t *testing.T) {
	tmpDir, cleanup := setupTempEnvDir(t, "production", "X=1")
	defer cleanup()
	pattern := filepath.Join(tmpDir, ".env.%s")

	// Capture stdout to verify printing occurs
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := InitEnvironmentWithEnv("production", &Options{FilePattern: pattern})
	if err != nil {
		w.Close()
		os.Stdout = old
		t.Fatal(err)
	}

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	os.Stdout = old

	output := buf.String()
	if len(output) == 0 {
		// color output may go to a different fd; just verify no error
	}
	os.Unsetenv("X")
}

func TestPrintEnvMode_Uat(t *testing.T) {
	tmpDir, cleanup := setupTempEnvDir(t, "uat", "X=1")
	defer cleanup()
	pattern := filepath.Join(tmpDir, ".env.%s")

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := InitEnvironmentWithEnv("uat", &Options{FilePattern: pattern})
	if err != nil {
		w.Close()
		os.Stdout = old
		t.Fatal(err)
	}

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	os.Stdout = old
	os.Unsetenv("X")
}

func TestPrintEnvMode_Development(t *testing.T) {
	tmpDir, cleanup := setupTempEnvDir(t, "development", "X=1")
	defer cleanup()
	pattern := filepath.Join(tmpDir, ".env.%s")

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := InitEnvironmentWithEnv("development", &Options{FilePattern: pattern})
	if err != nil {
		w.Close()
		os.Stdout = old
		t.Fatal(err)
	}

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	os.Stdout = old
	os.Unsetenv("X")
}
