package parse

import (
	"bytes"
	"fmt"
	"testing"
)

func TestLexerScanning(t *testing.T) {
	content := `model Page {
  fields {
    name: string
    content: string
  }
  pagination {
    per_page: 10
    max_per_page: 50
  }
}
`
	lexer := NewLexer("album.sb", bytes.NewBuffer([]byte(content)))
	scanner := lexer.Scan()
	for {
		lexeme := <-scanner
		fmt.Println(lexeme.Token, "('"+lexeme.Value+"')", lexeme.Pos)
		if lexeme.Token == EOF {
			break
		}
	}
}
