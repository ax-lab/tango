package common

import (
	"flag"
	"os"
)

func GetOutputDir(name, description string) string {
	outputDirPtr := flag.String(name, "", description)
	flag.Parse()

	outputDir := *outputDirPtr
	if outputDir == "" {
		ExitWithError("invalid arguments: missing output directory")
	}

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		ExitWithError("failed to create output directory: %v", err)
	}

	return outputDir
}
