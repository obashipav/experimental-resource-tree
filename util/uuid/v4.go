package uuid

import "github.com/google/uuid"

func NewID() string {
	return uuid.New().String()
}
