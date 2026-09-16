package shared

import (
	"fmt"

	"github.com/google/uuid"
)

func RandomKey(key string) (string, error) {
	raw, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	randomKey := fmt.Sprintf("%s/%s", key, raw.String())

	return randomKey, nil
}
