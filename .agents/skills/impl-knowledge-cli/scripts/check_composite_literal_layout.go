package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

func main() {
	fileSet := token.NewFileSet()
	failed := false
	if err := filepath.WalkDir(".", func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || filepath.Ext(path) != ".go" || filepath.IsLocal(path) && filepath.Dir(path) == ".agents/skills/impl-knowledge-cli/scripts" {
			return nil
		}
		file, err := parser.ParseFile(fileSet, path, nil, 0)
		if err != nil {
			return err
		}
		ast.Inspect(file, func(node ast.Node) bool {
			literal, ok := node.(*ast.CompositeLit)
			if !ok {
				return true
			}
			lines := map[int]bool{}
			for _, element := range literal.Elts {
				field, ok := element.(*ast.KeyValueExpr)
				if !ok {
					continue
				}
				line := fileSet.Position(field.Pos()).Line
				if lines[line] {
					fmt.Printf("%s:%d: 複合リテラルの各フィールドは改行してください\n", path, line)
					failed = true
				}
				lines[line] = true
			}

			return true
		})

		return nil
	}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if failed {
		os.Exit(1)
	}
}
