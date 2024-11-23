package lexer

import (
	f "fmt"
	"go/scanner"
	"go/token"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func ListGoFiles(path string) ([]string, error) {
	files := []string{}

	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(d.Name(), ".go") {
			return nil
		}

		files = append(files, path)
		return nil
	})

	return files, nil
}

// Count tokens in a go file in a special way: package declarations and
// imports are intentionally not counted. So are semicolons that are actually
// newlines. This way: refactoring your go file into multiple files does adds
// no additional tokens due to the extra package declaration and extra
// (duplicate) imports. Comments are also ignored.
func CountTokens(path string) (int, error) {

	count := 0

	fileset := token.NewFileSet()
	contents, err := os.ReadFile(path)

	if err != nil {
		return 0, f.Errorf("cannot read file: %w", err)
	}

	file := fileset.AddFile(path, fileset.Base(), len(contents))

	scan := scanner.Scanner{}
	scan.Init(file, contents, nil, scanner.ScanComments)

	insideImportSection := false

	for {
		_, tok, literal := scan.Scan()

		slog.Debug("counting tokens", "token", tok.String(), "literal", literal)

		if tok == token.EOF {
			break
		}

		// Don't count comments
		if tok == token.COMMENT {
			continue
		}

		// Don't count package declarations
		if tok == token.PACKAGE {
			_, pkgnameTok, _ := scan.Scan()

			if !pkgnameTok.IsLiteral() {
				slog.Warn("found non-literal when matching package token", "tokentype", pkgnameTok.String())
				continue
			}

			_, endlineTok, endlineLiteral := scan.Scan()
			if endlineTok != token.SEMICOLON || endlineLiteral != "\n" {
				slog.Warn("found non-semicolon when matching package token", "tokentype", endlineTok.String(), "literal", endlineLiteral)
				continue
			}

			// already done parsing this line
			continue
		}

		// Don't count newlines
		if tok == token.SEMICOLON && literal == "\n" {
			continue
		}

		// Don't count imports. Imports can have two formats:
		// import "package"
		// or
		// import (
		//   "package1"
		//   "package2"
		// )
		// This is implemented by ignoring up to and including the right
		// parenthesis in case of the multiline import. In case of an import
		// not followed by an opening parenthesis (the single line import) the
		// import statement is ignored.
		//
		// Need to take care though, that named imports of the form
		// import alias "actaulname" are correctly dealt with
		if tok == token.IMPORT {
			_, nextTok, _ := scan.Scan()

			if nextTok.IsLiteral() {
				// Skip named imports of the form
				// import f "fmt"
				_, nextTok, _ = scan.Scan()
			}

			if nextTok == token.LPAREN {
				insideImportSection = true
			}

			continue
		} else if insideImportSection {
			if tok == token.RPAREN {
				insideImportSection = false
			}

			continue
		}

		count += 1

	}

	return count, nil
}
