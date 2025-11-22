package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMainRuns(t *testing.T) {
	origArgs := os.Args
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("unable to determine working directory: %v", err)
	}

	os.Args = []string{"example.test"}

	rootDir := filepath.Dir(origDir)
	if err := os.Chdir(rootDir); err != nil {
		t.Fatalf("unable to change directory for test: %v", err)
	}

	defer func() {
		os.Args = origArgs
		_ = os.Chdir(origDir)
	}()

	main()
}
