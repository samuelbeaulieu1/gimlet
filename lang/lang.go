package lang

type Lang int

const (
	FR Lang = iota
	EN
)

var lang = EN

func Set(l Lang) {
	lang = l
}

func Get() Lang {
	return lang
}

func (l Lang) Value() int {
	return int(l)
}
