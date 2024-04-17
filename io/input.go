package io

import (
	"encoding/json"
	"io"
	"os"
	"v0/scene"

	"go.uber.org/zap"
)

type Input struct {
	Camera scene.Camera

	AspectRatio float64
	ImageWidth  int

	Scene scene.Scene

	NumWorkers int

	// The sqrt of the number of subpixels
	AntiAliasCount int

	MaxDepth int
}

func UnmarshalInput(filePath string) *Input {
	fileErrFunc := func(msg string, err error) {
		Logger.Panic(msg,
			zap.String("filePath", filePath),
			zap.Error(err),
		)
	}
	f, err := os.Open(filePath)
	if err != nil {
		fileErrFunc("input open err", err)
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		fileErrFunc("input read err", err)
	}

	var input Input
	if err := json.Unmarshal(bytes, &input); err != nil {
		fileErrFunc("input unmarshal err", err)
	}

	if input.NumWorkers <= 0 {
		Logger.Panic("invalid 'NumWorkers', must be > 0", zap.Int("NumWorkers", input.NumWorkers))
	}
	if input.MaxDepth <= 0 {
		Logger.Panic("invalid 'MaxDepth', must be > 0", zap.Int("MaxDepth", input.NumWorkers))
	}
	if input.AntiAliasCount < 1 {
		Logger.Panic("invalid 'AntiAliasCount', must be > 0", zap.Int("AntiAliasCount", input.AntiAliasCount))
	}

	return &input
}
