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

// SettingsFieldDoesNotExist is an error when a field is specified via the DefaultsMap that does not exist
// in the out struct value that is provided to settings.Gather
func SettingsFieldDoesNotExist(overrideType string, fieldName string) SettingsError {
	return SettingsError{
		Message: fmt.Sprintf("field specified in override (%s) does not exist in the target out struct: %s", overrideType, fieldName),
	}
}

// SettingsFieldTypeMismatch is raised in the event there is a mismatch between types when trying to override a specific value
func SettingsFieldTypeMismatch(fieldName string, expectedType reflect.Kind, receivedType reflect.Kind) SettingsError {
	return SettingsError{
		Message: fmt.Sprintf("type mismatch for field %s: expected %v but value is %v", fieldName, expectedType, receivedType),
	}
}

// SettingsFieldSetError is raised when a field's value is not actually settable
func SettingsFieldSetError(fieldName string, t reflect.Kind, m ...error) SettingsError {
	if len(m) == 0 {
		return SettingsError{
			Message: fmt.Sprintf("unable to set the value of a field in settings: %s (type: %v)", fieldName, t),
		}
	}

	return SettingsError{
		Message: fmt.Sprintf(
			"unable to set the value of a field in settings: %s (type: %v): %s",
			fieldName,
			t,
			m[0].Error()),
	}
}

// SettingsFileParseError occurs when a specified settings file can't be properly unmarshalled
func SettingsFileParseError(path string, desc string) SettingsError {
	return SettingsError{
		Message: fmt.Sprintf("unable to parse settings file (%s): %s", path, desc),
	}
}

// SettingsFileReadError occurs when a specified settings file is not readable
func SettingsFileReadError(path string, desc string) SettingsError {
	return SettingsError{
		Message: fmt.Sprintf("unable to read settings file (%s): %s", path, desc),
	}
}

// SettingsFileTypeError occurs when a format is requested that the settings package does not support
func SettingsFileTypeError(path string, ext string) SettingsError {
	return SettingsError{
		Message: fmt.Sprintf("unrecognized settings file extension (%s): %s", path, ext),
	}
}

// SettingsOutCannotBeNil occurs when the out field in the settings struct is set to nil, intentionally or otherwise
func SettingsOutCannotBeNil() SettingsError {
	return SettingsError{
		Message: "out cannot be nil",
	}
}

// SettingsTypeDiscoveryError occurs when the out value provided to settings.Gather is not a struct
func SettingsTypeDiscoveryError(t reflect.Kind) SettingsError {
	return SettingsError{
		Message: fmt.Sprintf("unable to detect fields for non-struct type: %v", t),
	}
}
