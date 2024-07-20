package utility

import "github.com/gofrs/uuid"

func GenerateUUID() string {
	id, _ := uuid.NewV7()
	return id.String()
}
