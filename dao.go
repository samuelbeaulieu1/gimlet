package gimlet

import (
	"math/rand"
	"time"
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

func CreateNewID(dao Dao, model Model, n int) string {
	var charset = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_")
	id := make([]byte, n)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		for i := range id {
			id[i] = charset[r.Intn(len(charset))]
		}
		if !dao.ExistsByID(string(id), model) {
			return string(id)
		} else {
			id = make([]byte, n)
		}
	}
}
