package testutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	wd := BasePath()

	env := filepath.Join(wd, ".env")
	if err := godotenv.Load(env); err != nil {
		panic(err)
	}

	environment := os.Getenv("ENVIRONMENT")
	if environment == "local" {
		envLocal := filepath.Join(wd, ".env.local")
		if _, err := os.Stat(envLocal); !errors.Is(err, os.ErrNotExist) {
			if err := godotenv.Overload(envLocal); err != nil {
				fmt.Println("Error: .env.local is not well formatted, cannot overload values.")
				panic(err)
			}
		}
	}
}
