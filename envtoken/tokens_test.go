package envtoken

import (
	"fmt"
	"os"
	"testing"
)

func TestEnvironmentInvalidNotSet(t *testing.T) {
	envTokens := []*EnvToken{
		NewEnvToken("_TEST_NOT_SET_1", "", true),
		NewEnvToken("_TEST_NOT_SET_2", "", true),
	}
	unsetEnvTokens(envTokens)

	_, errCollection := NewEnvironment(envTokens)
	if errCollection == nil {
		t.Fatalf("Error expected, none returned")
	}

	expected := len(envTokens)
	assertUnsetKeyErrors(t, errCollection, expected)

	expected = 0
	assertEmptyValueErrors(t, errCollection, expected)
}

func TestEnvironmentInvalidEmptyValue(t *testing.T) {
	envTokens := []*EnvToken{
		NewEnvToken("_TEST_NOT_SET_1", "", true),
		NewEnvToken("_TEST_NOT_SET_2", "", true),
	}
	expected := 0
	tokens := []*EnvToken{}
	for i := range envTokens {
		token := envTokens[i]
		if token.Required {
			tokens = append(tokens, token)
			err := os.Setenv(token.EnvKey, "")
			if err != nil {
				t.Fatalf("Error setting up environment, %v", err)
			}
			expected++
		}
	}

	defer unsetEnvTokens(envTokens)
	_, errCollection := NewEnvironment(envTokens)
	if errCollection == nil {
		t.Fatalf("Error expected, none returned")
	}

	assertEmptyValueErrors(t, errCollection, expected)
}

func TestEnvironment(t *testing.T) {
	envTokens := []*EnvToken{
		NewEnvToken("_TEST_NOT_SET_1", "", true),
		NewEnvToken("_TEST_NOT_SET_2", "", true),
		NewEnvToken("_TEST_NOT_SET_3", "", false),
	}
	defer unsetEnvTokens(envTokens)

	tokens := []*EnvToken{}
	for i := range envTokens {
		token := envTokens[i]
		if token.Required {
			tokens = append(tokens, token)
			err := os.Setenv(token.EnvKey, fmt.Sprintf("%v", token.EnvKey))
			if err != nil {
				t.Fatalf("Error setting up environment, %v", err)
			}
		}
	}

	defer unsetEnvTokens(tokens)
	env, errCollection := NewEnvironment(tokens)
	if errCollection != nil {
		t.Fatalf("Error not expected, %v", errCollection.GetError())
	}

	for i := range tokens {
		token := tokens[i]
		actual := *env.GetValue(token.EnvKey)
		expected := token.EnvKey
		if actual != expected {
			t.Fatalf("%v has unexpected value, expected %v, actual %v", token, actual, expected)
		}

	}
}
func TestEnvironmentDefaults(t *testing.T) {
	envTokens := []*EnvToken{
		NewEnvToken("_TEST_NOT_SET_1", "1", false),
		NewEnvToken("_TEST_NOT_SET_2", "2", false),
		NewEnvToken("_TEST_NOT_SET_3", "4", false),
	}

	unsetEnvTokens(envTokens)

	env, errCollection := NewEnvironment(envTokens)
	if errCollection != nil {
		t.Fatalf("Error not expected, %v", errCollection.GetError())
	}

	for i := range envTokens {
		token := envTokens[i]
		actual := *env.GetValue(token.EnvKey)
		expected := *token.DefaultValue
		if actual != expected {
			t.Fatalf("%v has unexpected value, expected %v, actual %v", token, actual, expected)
		}
	}

	token := envTokens[0]
	err := os.Setenv(token.EnvKey, fmt.Sprintf("%v", token.EnvKey))
	if err != nil {
		t.Fatalf(err.Error())
	}
	env, errCollection = NewEnvironment(envTokens)
	if errCollection != nil {
		t.Fatalf("Error not expected, %v", errCollection.GetError())
	}

	defer unsetEnvTokens(envTokens)

}

func TestBoolFromEnvFalse(t *testing.T) {
	expected := false
	key := "_BOOL_FALSE"
	defer os.Unsetenv(key)

	os.Unsetenv(key)
	value := ""
	assertBoolFromEnvValue(t, key, value, expected)

	value = "no"
	os.Setenv(key, value)
	assertBoolFromEnvValue(t, key, value, expected)
}

func assertBoolFromEnvValues(t *testing.T, key string, expected bool, values []string) {
	t.Helper()
	for i := range values {
		value := values[i]
		os.Setenv(key, value)
		assertBoolFromEnvValue(t, key, value, expected)
	}
}

func assertBoolFromEnvValue(t *testing.T, key, value string, expected bool) {
	t.Helper()
	actual := BoolFromEnv(key)
	if actual != expected {
		t.Fatalf("key %v=%v expected %v, actual %v", key, value, expected, actual)
	}
}

func assertEmptyValueErrors(t *testing.T, errCollection *EnvErrorCollection, expected int) {
	t.Helper()
	actual := len(errCollection.emptyValueErrors)
	if expected != actual {
		t.Fatalf("Expected %v empty value errors found %v", expected, actual)
	}
}

func assertUnsetKeyErrors(t *testing.T, errCollection *EnvErrorCollection, expected int) {
	t.Helper()
	actual := len(errCollection.unsetErrors)
	if expected != actual {
		t.Fatalf("Expected %v unset key errors not found %v", expected, actual)
	}
}

func unsetEnvTokens(tokens []*EnvToken) {
	for i := range tokens {
		unsetEnv(tokens[i].EnvKey)
	}
}

func unsetEnv(key string) error {
	err := os.Unsetenv(key)
	if err != nil {
		fmt.Printf("failed unsetting env %v", err)
	}
	return err
}
