package wuuid

import "github.com/google/uuid"

func NewUUID() uuid.UUID {
    return uuid.Must(uuid.NewV7())
}
