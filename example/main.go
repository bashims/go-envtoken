package main

import (
	"fmt"
	"os"

	"github.com/bashims/go-envtoken/envtoken"
)

func main() {
	envTokens := []*envtoken.EnvToken{
		envtoken.NewEnvToken("KEY1", "default", false),
		envtoken.NewEnvToken("KEY2", "", true),
	}

	env, errCollection := envtoken.NewEnvironment(envTokens)
	if errCollection != nil {
		fmt.Println(errCollection.GetError())
	}

	os.Setenv("KEY2", "value2")

	env, errCollection = envtoken.NewEnvironment(envTokens)
	fmt.Printf("KEY1=%s\n", *env.GetValue("KEY1"))
	fmt.Printf("KEY2=%s\n", *env.GetValue("KEY2"))
}
