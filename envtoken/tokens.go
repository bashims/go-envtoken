package envtoken

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	reBoolean = regexp.MustCompile("^((?i)true|(?i)yes|1)$")
)

//EnvToken
type EnvToken struct {
	EnvKey       string
	Required     bool
	DefaultValue *string
	Value        *string
}

func (e *EnvToken) SetValue(v *string) {
	e.Value = v
}

//
func NewEnvToken(envKey string, defaultValue string, required bool) *EnvToken {
	return &EnvToken{
		EnvKey:       envKey,
		DefaultValue: &defaultValue,
		Required:     required,
	}
}

// Environment configuration
type Environment struct {
	m map[string]*EnvToken
}

// Add key with value to Environment
func (e *Environment) Add(t *EnvToken) {
	e.m[t.EnvKey] = t
}

// Get *EnvToken for environment key.
func (e *Environment) Get(a string) *EnvToken {
	if t, ok := e.m[a]; ok {
		return t
	}
	return nil
}

// GetValue return either the set value or default value for environment key.
func (e *Environment) GetValue(k string) *string {
	t := e.Get(k)
	if t == nil {
		return nil
	}

	if *t.Value != "" {
		return t.Value
	}
	if *t.DefaultValue != "" {
		return t.DefaultValue
	}
	return nil
}

// GetBoolean return
func (e *Environment) GetBoolean(k string) bool {
	v := e.GetValue(k)
	if v == nil {
		return false
	}
	return GetBoolean(*v)
}

// EnvKeyNotSetError denotes an missing environment key
type EnvKeyNotSetError struct {
	key string
}

// Error returns the error string.
func (e *EnvKeyNotSetError) Error() string {
	return fmt.Sprintf("key %v, not set", e.key)
}

// NewEnvKeyNotSetError return a ErrEnvKeyNotSet for environement key.
func NewEnvKeyNotSetError(key string) *EnvKeyNotSetError {
	return &EnvKeyNotSetError{key: key}
}

// EnvEmptyValueError denotes an environment key with an empty value
type EnvEmptyValueError struct {
	key string
}

// Error returns the error string.
func (e *EnvEmptyValueError) Error() string {
	return fmt.Sprintf("key %v, value is empty", e.key)
}

// NewErrorEnvEmptyValue return a ErrorEnvEmptyValue for environment key.
func NewErrorEnvEmptyValue(key string) *EnvEmptyValueError {
	return &EnvEmptyValueError{key: key}
}

// EnvErrorCollection provides a way to collet EnvKeyNotSetError and EnvEmptyValueError
type EnvErrorCollection struct {
	unsetErrors      []*EnvKeyNotSetError
	emptyValueErrors []*EnvEmptyValueError
}

// AddKeyNotSet add a EnvKeytNotSetError to EnvErrorCollection, return the error
func (e *EnvErrorCollection) AddKeyNotSet(key string) *EnvKeyNotSetError {
	err := NewEnvKeyNotSetError(key)
	e.unsetErrors = append(e.unsetErrors, err)
	return err
}

// AddKeyEmptyValue add a EnvEmptyValueError to EnvErrorCollection, return the error
func (e *EnvErrorCollection) AddKeyEmptyValue(key string) *EnvEmptyValueError {
	err := NewErrorEnvEmptyValue(key)
	e.emptyValueErrors = append(e.emptyValueErrors, err)
	return err
}

// GetError return the collection of all errors as error.
func (e *EnvErrorCollection) GetError() error {
	errs := []string{}
	for i := range e.emptyValueErrors {
		errs = append(errs, e.emptyValueErrors[i].Error())

	}
	for i := range e.unsetErrors {
		errs = append(errs, e.unsetErrors[i].Error())

	}
	count := len(errs)
	if count > 0 {
		return fmt.Errorf("Invalid environment, %d error(s):\n%v", count, strings.Join(errs, "\n"))
	}
	return nil
}

// NewEnvErrorCollection return a new EnvErrorCollection
func NewEnvErrorCollection() EnvErrorCollection {
	return EnvErrorCollection{
		unsetErrors:      []*EnvKeyNotSetError{},
		emptyValueErrors: []*EnvEmptyValueError{},
	}

}

// NewEnvironment returns a new Environment
func NewEnvironment(envTokens []*EnvToken) (*Environment, *EnvErrorCollection) {
	env := &Environment{
		m: make(map[string]*EnvToken),
	}

	errCollection := NewEnvErrorCollection()
	for i := range envTokens {
		t := envTokens[i]
		v, exists := os.LookupEnv(t.EnvKey)
		if t.Required {
			if !exists {
				_ = errCollection.AddKeyNotSet(t.EnvKey)
				continue
			} else if v == "" {
				_ = errCollection.AddKeyEmptyValue(t.EnvKey)
				continue
			}
		}
		t.SetValue(&v)
		env.Add(t)
	}

	if errCollection.GetError() == nil {
		return env, nil
	}
	return env, &errCollection
}

//BoolFromEnv return boolean value based on the the environment key's value.
func BoolFromEnv(key string) bool {
	return GetBoolean(os.Getenv(key))
}

//GetBoolean return boolean value based on value.
func GetBoolean(value string) bool {
	return reBoolean.MatchString(value)
}
