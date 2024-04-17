package io

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

type Environment struct {
	InputFilePath  string `env:"INPUT_FILE_PATH,default=input.json"`
	OutputFilePath string `env:"OUTPUT_FILE_PATH,default=output.ppm"`
}

var Env = func() Environment {
	ctx := context.Background()
	var env Environment
	if err := envconfig.Process(ctx, &env); err != nil {
		panic(err)
	}

	return env
}()
