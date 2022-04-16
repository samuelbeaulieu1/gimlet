package gimlet

type Model interface {
	ToDTO() DTO
	TableName() string
}
