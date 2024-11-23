package util

import "github.com/google/uuid"

func GenerateVideoUUID() string {
	newUUID := uuid.New()
	return "VIDEO" + newUUID.String()
}
