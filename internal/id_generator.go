package internal

type IdGenerator interface {
	Generate() (string, error)
}
