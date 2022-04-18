package gimlet

type Authenticator interface {
	GetAuth() any
	SetAuth(any)
	ValidAuth() error
}
