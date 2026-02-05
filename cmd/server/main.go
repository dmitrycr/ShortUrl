package main

import (
	"fmt"

	"github.com/dmitrycr/ShortUrl/pkg/generator"
)

func main() {
	g := generator.NewGenerator(7)

	fmt.Println(g.Generate())
}
