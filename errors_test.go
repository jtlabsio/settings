package settings

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestSettingsFieldSetError(t *testing.T) {
	errNoMessage := SettingsFieldSetError("Name", reflect.String)
	if !strings.Contains(errNoMessage.Error(), "Name") || strings.Contains(errNoMessage.Error(), "boom") {
		t.Fatalf("SettingsFieldSetError() without underlying error = %v", errNoMessage)
	}

	underlying := errors.New("boom")
	errWithMessage := SettingsFieldSetError("Name", reflect.String, underlying)
	if !strings.Contains(errWithMessage.Error(), underlying.Error()) {
		t.Fatalf("SettingsFieldSetError() with underlying error missing detail: %v", errWithMessage)
	}
}

func TestSettingsFileReadError(t *testing.T) {
	err := SettingsFileReadError("/tmp/config.yml", "permission denied")
	if !strings.Contains(err.Error(), "permission denied") {
		t.Fatalf("SettingsFileReadError() = %v", err)
	}
}
