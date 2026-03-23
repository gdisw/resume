package identifier

import (
	"github.com/google/uuid"
	"github.com/segmentio/ksuid"
)

func Generate() string {
	return ksuid.New().String()
}

func GenerateUUID() string {
	return uuid.New().String()
}
