package uuid

import (
	uuid "github.com/satori/go.uuid"
)

func GetUuid() (uuidString string) {
	uuidString = uuid.Must(uuid.NewV4(), nil).String()
	return
}
