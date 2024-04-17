package main

import (
	"v0/io"
	"v0/tracer"
)

func main() {
	defer io.InitLogger()()

	input := io.UnmarshalInput(io.Env.InputFilePath)

	tracer.RayTrace(input)
}
