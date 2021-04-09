package settings

import (
	"fmt"
	"reflect"
)

type SettingsError struct {
	Message string
}

func (e SettingsError) Error() string {
	return e.Message
}

func SettingsFieldDoesNotExist(overrideType string, fieldName string) SettingsError {
	return SettingsError{
		Message: fmt.Sprintf("field specified in override (%s) does not exist in the target out struct: %s", overrideType, fieldName),
	}
}

func SettingsFieldSetError(fieldName string, t reflect.Kind) SettingsError {
	return SettingsError{
		Message: fmt.Sprintf("unable to set the value of a field in settings: %s (type: %v)", fieldName, t),
	}
}

func SettingsFileParseError(path string, desc string) SettingsError {
	return SettingsError{
		Message: fmt.Sprintf("unable to parse settings file (%s): %s", path, desc),
	}
}

func SettingsFileReadError(path string, desc string) SettingsError {
	return SettingsError{
		Message: fmt.Sprintf("unable to read settings file (%s): %s", path, desc),
	}
}

func SettingsFileTypeError(path string, ext string) SettingsError {
	return SettingsError{
		Message: fmt.Sprintf("unrecognized settings file extension (%s): %s", path, ext),
	}
}

func SettingsTypeDiscoveryError(t reflect.Kind) SettingsError {
	return SettingsError{
		Message: fmt.Sprintf("unable to detect fields for non-struct type: %v", t),
	}
}
