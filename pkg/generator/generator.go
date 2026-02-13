package generator

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

const (
	base62Chars   = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	defaultLength = 6 // Длина кода по умолчанию
	//maxRetries    = 3
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

	var code strings.Builder
	code.Grow(g.length)

	charsetLength := big.NewInt(int64(len(base62Chars)))

	for i := 0; i < g.length; i++ {
		// Генерируем криптографически стойкое случайное число
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}

		code.WriteByte(base62Chars[randomIndex.Int64()])
	}

	return code.String(), nil
}

// EncodeID конвертирует числовой ID в Base62 строку
// Полезно для использования с автоинкрементом из БД
func (g *Generator) EncodeID(id int64) string {
	if id == 0 {
		return string(base62Chars[0])
	}

	var result strings.Builder
	base := int64(len(base62Chars))

	for id > 0 {
		remainder := id % base
		result.WriteByte(base62Chars[remainder])
		id = id / base
	}

	// Переворачиваем строку (т.к. собирали в обратном порядке)
	return reverse(result.String())
}

// DecodeID конвертирует Base62 строку обратно в числовой ID
func (g *Generator) DecodeID(code string) (int64, error) {
	var id int64
	base := int64(len(base62Chars))

	for _, char := range code {
		// Находим позицию символа в алфавите
		pos := strings.IndexRune(base62Chars, char)
		if pos == -1 {
			return 0, fmt.Errorf("ERROR")
		}

		id = id*base + int64(pos)
	}

	return id, nil
}

// ValidateCode проверяет, что код содержит только допустимые символы
func (g *Generator) ValidateCode(code string) bool {
	if len(code) == 0 || len(code) > 20 {
		return false
	}

	for _, char := range code {
		if !strings.ContainsRune(base62Chars, char) {
			return false
		}
	}

	return true
}

// reverse переворачивает строку
func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
