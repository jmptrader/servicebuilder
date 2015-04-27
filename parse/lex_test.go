package parse

import (
	"bytes"
	"fmt"
	"testing"
)

func TestLexerScanning(t *testing.T) {
	content := `Page {
  fields {
    name: string
    content: string
  }
  pagination {
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
