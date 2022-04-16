package actions

type Action int

const (
	ReadAction Action = iota
	CreateAction
	UpdateAction
	DeleteAction
)
