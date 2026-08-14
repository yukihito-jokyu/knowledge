#!/usr/bin/env python3
"""Check Go composite literal element layout with the Go parser."""

from __future__ import annotations

import subprocess
import sys
import tempfile
from pathlib import Path


GO_CHECKER = r'''package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: checker <repository-root>")
		os.Exit(2)
	}
	root := os.Args[1]
	fileSet := token.NewFileSet()
	failed := false
	reported := map[string]bool{}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if path != root && (entry.Name() == ".git" || entry.Name() == ".agents") {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}
		file, parseErr := parser.ParseFile(fileSet, path, nil, 0)
		if parseErr != nil {
			return parseErr
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
					relative, relErr := filepath.Rel(root, path)
					if relErr != nil {
						relative = path
					}
					key := fmt.Sprintf("%s:%d", relative, line)
					if !reported[key] {
						fmt.Printf("%s: 複合リテラルの各フィールドは改行してください\n", key)
						reported[key] = true
					}
					failed = true
				}
				lines[line] = true
			}
			return true
		})
		return nil
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if failed {
		os.Exit(1)
	}
}
'''


def main() -> int:
    repo_result = subprocess.run(
        ("git", "rev-parse", "--show-toplevel"),
        check=False,
        text=True,
        capture_output=True,
    )
    if repo_result.returncode != 0:
        sys.stderr.write(repo_result.stderr)
        return repo_result.returncode
    repo = Path(repo_result.stdout.strip()).resolve()

    with tempfile.TemporaryDirectory(prefix="knowledge-literal-check-") as temp_dir:
        checker = Path(temp_dir) / "main.go"
        checker.write_text(GO_CHECKER, encoding="utf-8")
        result = subprocess.run(("go", "run", checker.as_posix(), repo.as_posix()), cwd=repo)
        return result.returncode


if __name__ == "__main__":
    raise SystemExit(main())
