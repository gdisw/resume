package testutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSuiteBehaviour(t *testing.T) {
	var checked int
	require.NoError(t, filepath.Walk("..", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if m, err := filepath.Match("*_test.go", info.Name()); err != nil || !m {
			return err
		}

		checked++

		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		text := string(b)
		if !strings.Contains(text, "TestSuite struct") {
			return nil
		}

		if strings.Contains(text, "TestSuite(t *testing.T)") {
			return nil
		}

		assert.FailNow(t, fmt.Sprintf("TestSuite struct in %s does not have TestSuite(t *testing.T) method", path))
		return nil
	}))

	assert.Greater(t, checked, 0)
}
