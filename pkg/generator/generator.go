package generator

import (
	"strconv"
	"time"
)

const (
	Base62        = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	defaultLength = 6 // Длина кода по умолчанию
	maxRetries    = 3
)

type Generator struct {
	length int
}

type CodeGenerator interface {
	GenerateUnique() (string, error)
}

func NewGenerator(length int) *Generator {
	if length <= 0 {
		length = defaultLength
	}

	return &Generator{length: length}
}

func (g *Generator) Generate() (string, error) {
	if g.length <= 0 {
		return "", ErrInvalidLength
	}

	code := "shortUrl" + strconv.Itoa(time.Now().Nanosecond())

	return code, nil
}
