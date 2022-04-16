package gimlet

type Entity[M any] interface {
	Delete(id string) error
	Update(id string, request *M) error
	Create(request *M) (*M, error)
	Get(id string) (*M, error)
	GetAll() (*[]M, error)
	Exists(id string) bool
}
