package testutil

import (
	"os"
	"path/filepath"
)

func BasePath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for filepath.Base(wd) != "resume" {
		wd = filepath.Dir(wd)
	}

	return wd
}
