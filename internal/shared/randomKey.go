package shared

import (
	"fmt"

	"github.com/google/uuid"
)

func RandomKey(prefix string) (string, error) {
	raw, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	randomKey := fmt.Sprintf("%s/%s", prefix, raw.String())

	return randomKey, nil
}
