package nanoid

import "github.com/jaevor/go-nanoid"

type Generator struct {
	alphabet string
	length   int
}

func NewNanoIdGenerator(alphabet string, length int) *Generator {
	return &Generator{
		alphabet: alphabet,
		length:   length,
	}
}

func (generator Generator) Generate() (string, error) {
	gen, err := nanoid.CustomASCII(generator.alphabet, generator.length)

	if err != nil {
		return "", err
	}

	return gen(), nil
}
