package gimlet

import (
	"encoding/base64"

	"github.com/google/uuid"
)

const (
	DefaultIDLength = 12
)

type Dao interface {
	Delete(id string, model Model) error
	Create(model Model) error
	Update(model Model, update Model) error
	Get(id string, model Model) error
	ExistsByID(id string, model Model) bool
}

func CreateNewID() (string, error) {
	id := uuid.New()
	uuidBytes, err := id.MarshalBinary()
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString([]byte(uuidBytes)), nil
}
