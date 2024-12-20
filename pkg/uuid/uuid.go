package uuid

import "github.com/google/uuid"

type Uuid struct {
}

func New() *Uuid {
	return &Uuid{}
}

func (u *Uuid) NewString() string {
	return uuid.NewString()
}

func ParseUUID(uuidStr string) (uuid.UUID, error) {
	return uuid.Parse(uuidStr)
}
