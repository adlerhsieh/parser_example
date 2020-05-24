package main

import (
	"github.com/adlerhsieh/parser_example/extraline"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(extraline.Analyzer)
}
