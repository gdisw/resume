package identifier

import (
	"fmt"

	"github.com/segmentio/ksuid"
)

func Mint(prefix string) string {
	id := ksuid.New().String()
	return fmt.Sprintf("%s_%s", prefix, id)
}
