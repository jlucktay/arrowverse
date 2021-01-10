package main

import (
	"fmt"
	"os"

	"go.uber.org/zap"

	"go.jlucktay.dev/arrowverse/cmd"
)

func main() {
	logger, errNP := zap.NewProduction()
	if errNP != nil {
		fmt.Fprintf(os.Stderr, "could not create new logger: %v", errNP)
		return
	}
	defer logger.Sync() //nolint:errcheck

	if errExec := cmd.Execute(); errExec != nil {
		logger.Error("failed in execute",
			// Structured context as strongly typed Field values.
			zap.Error(errExec),
		)
	}
}
